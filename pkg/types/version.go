package types

type Version struct {
	Codes []Number `json:"codes" protobuf:"bytes,1,opt,name=codes"`
}

type Number struct {
	Name  string `json:"name" protobuf:"bytes,1,opt,name=name"`
	Value int32  `json:"value" protobuf:"varint,2,opt,name=value"`
}
