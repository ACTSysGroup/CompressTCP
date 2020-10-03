package common
type DatasetMeta struct {
	Key    string `yaml:"key"`
	Source string `yaml:"source"`
	File   string `yaml:"file"`
	Desc   string `yaml:"desc"`
	Size   string `yaml:"size"`
}
