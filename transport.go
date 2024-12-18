package mediasoupgo

type Transport struct{}

func (t *Transport) Dump() {}

func (t *Transport) Produce() {}

func (t *Transport) ProduceData() {}

func (t *Transport) Consume() {}

func (t *Transport) ConsumeData() {}
