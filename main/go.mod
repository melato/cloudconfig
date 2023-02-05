module main

go 1.19

replace (
	melato.org/cloudconfig => ../
	melato.org/yaml => ../../yaml

)

require (
	melato.org/cloudconfig v0.0.0-00010101000000-000000000000
	melato.org/command v1.0.0
	melato.org/yaml v0.0.0-00010101000000-000000000000
)

require (
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
