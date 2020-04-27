package slice

import "reflect"

// Reverse sort the given slice by reversing the order
func Reverse(a interface{}) {
	switch reflect.TypeOf(a).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(a)

		for i := s.Len()/2 - 1; i >= 0; i-- {
			opp := s.Len() - 1 - i
			lastIndex := s.Index(opp)
			tmpVal := lastIndex.Interface()
			firstIndex := s.Index(i)
			lastIndex.Set(firstIndex)
			firstIndex.Set(reflect.ValueOf(tmpVal))
		}
	}
}
