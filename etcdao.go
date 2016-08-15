package etcdao

import 
	"golang.org/x/net/context"
	"github.com/coreos/etcd/client"
	"reflect"
	"strconv"
	"errors"
	"strings"
	"time"
)

var ErrBadFormat = errors.New("ErrBadFormat")

func processNode(node *client.Node,v interface {},tag reflect.StructTag) error{
	var err error
	switch v := v.(type) {
	case *int:
		*v,err = strconv.Atoi(node.Value)
		if err != nil {
			log.Warnf("processNode: node.Value is not int '%v'",node.Value)
		}
	case *string:
		*v = node.Value
	case *time.Time:
		format := tag.Get("format")
		if format == "" {
			format = "2006-01-02"
		}
		*v , err = time.Parse(format,node.Value)
		if err != nil {
			log.WithError(err).Warnf("processNode: node.Value is not time '%v' for format '%s'",node.Value,format)
		}
	case *bool:
		if node.Value != "" && node.Value != "0" && node.Value != "false"  {
			*v = true
		} else {
			*v = false
		}
	default:
		tp := reflect.TypeOf(v).Elem()
		switch tp.Kind() {
		case reflect.Map:
			if ! node.Dir {
				log.Warnf("processNode: type %v needs dir in etcd , not scalar",reflect.TypeOf(v))
				return ErrBadFormat
			}

			mapElem := reflect.MakeMap(tp)

			//maxVal := 0
			for _,v := range node.Nodes {
				li := strings.LastIndex(v.Key,"/")
				fieldName := v.Key[li+1:]
				newValPtr := reflect.New(tp.Elem())

				err = processNode(v,newValPtr.Interface(),tag)
				if err != nil {
					return err
				}

				mapElem.SetMapIndex(reflect.ValueOf(fieldName),newValPtr.Elem())

			}


			reflect.ValueOf(v).Convert(reflect.TypeOf(v)).Elem().Set(mapElem)

		case reflect.Slice:
			if ! node.Dir {
				log.Warnf("processNode: type %v needs dir in etcd , not scalar",reflect.TypeOf(v))
				return ErrBadFormat
			}

			dirFields := make(map[int]*client.Node)
			maxVal := 0
			for _,v := range node.Nodes {
				li := strings.LastIndex(v.Key,"/")
				fieldName,err := strconv.Atoi(v.Key[li+1:])
				if err != nil {
					log.Warnf("processNode: dir element '%s' must be integer",v.Key[li+1:])
					return ErrBadFormat
				}
				dirFields[fieldName] = v
				if maxVal <  fieldName {
					maxVal = fieldName
				}
			}

			sliceElem := reflect.MakeSlice(tp,maxVal+1,maxVal+1)

			for i := 0 ; i <= maxVal ; i++ {
				if dirFields[i] != nil {
					curValue := sliceElem.Index(i)
					newValPtr := reflect.New(curValue.Type())

					err = processNode(dirFields[i],newValPtr.Interface(),tag)
					if err != nil {
						return err
					}
					curValue.Set(newValPtr.Elem())
				}
			}
			reflect.ValueOf(v).Convert(reflect.TypeOf(v)).Elem().Set(sliceElem)


		case reflect.Struct:
			if ! node.Dir {
				log.Warnf("processNode: type %v needs dir in etcd , not scalar",reflect.TypeOf(v))
				return ErrBadFormat
			}
			dirFields := make(map[string]*client.Node)
			for _,v := range node.Nodes {
				li := strings.LastIndex(v.Key,"/")
				dirFields[v.Key[li+1:]] = v
			}

			svals := reflect.ValueOf(v).Convert(reflect.TypeOf(v)).Elem()

			fc := tp.NumField()
			for i:= 0 ; i < fc ; i++ {
				fieldTp := tp.Field(i)
				//fieldElem := svals.Field(i)
				//log.Println(fieldElem)
				fieldName := fieldTp.Tag.Get("name")
				if fieldName == "" {
					fieldName = fieldTp.Name
				}
				if fieldNode,ok := dirFields[fieldName] ; ok {
					curValue := svals.Field(i)
					newValPtr := reflect.New(fieldTp.Type)

					err = processNode(fieldNode,newValPtr.Interface(),fieldTp.Tag)
					if err != nil {
						return err
					}
					curValue.Set(newValPtr.Elem()/*fieldVal.Elem()*/)
					//log.Warnf(" field '%s' type: %v curVal=%v",fieldTp.Name,fieldTp,newValPtr.Elem())
				}
			}
		default:
			log.Warnf("processNode: Unkonwn type %v KIND=%v",reflect.TypeOf(v),tp.Kind())
			return ErrBadFormat
		}
	}
	return nil
}

/*
	Reflexive object reader from etcd

	supports:
		int,string,time.Time,bool scalars
		structs with tagging info:
			name = if name of etcd field differs from struct name
			format = date conversion format for time.Time
		slices - reads array from etcd directory - awaits 0 , 1 , 2 ... etc as keys
		map

 */
func ReadObject(kapi,ctxt,path string, v interface {} ) error {
	resp, err := kapi.Get(ctxt, path, &client.GetOptions{Recursive:true})
	if err != nil {
		return err
	}

	return processNode(resp.Node,v,"")

}
