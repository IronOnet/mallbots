package es

import (
	"fmt"

	"github.com/irononet/mallbots/internal/registry"
)

type VersionSetter interface{
	setVersion(int)
}

func SetVersion(version int) registry.BuildOption{
	return func(v any) error{
		if agg, ok := v.(VersionSetter); ok{
			agg.setVersion(version)
			return nil
		}
		return fmt.Errorf("%T does not have the method setVersion(int)", v)
	}
}