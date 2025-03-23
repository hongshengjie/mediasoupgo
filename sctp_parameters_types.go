package mediasoupgo

// NumSctpStreams represents the number of outgoing and incoming SCTP streams.
type NumSctpStreams struct {
    // Initially requested number of outgoing SCTP streams (OS).
    OS int
    // Maximum number of incoming SCTP streams (MIS).
    MIS int
}

// SctpCapabilities represents SCTP capabilities.
type SctpCapabilities struct {
    // Number of SCTP streams.
    NumStreams NumSctpStreams
}

// SctpParameters represents SCTP parameters.
type SctpParameters struct {
    // Must always equal 5000.
    Port int
    // Initially requested number of outgoing SCTP streams.
    OS int
    // Maximum number of incoming SCTP streams.
    MIS int
    // Maximum allowed size for SCTP messages.
    MaxMessageSize int
}

// SctpStreamParameters represents SCTP stream parameters describing the reliability of a certain SCTP stream.
// If Ordered is true, MaxPacketLifeTime and MaxRetransmits must be false.
// If Ordered is false, only one of MaxPacketLifeTime or MaxRetransmits can be true.
type SctpStreamParameters struct {
    // SCTP stream id.
    StreamId int
    // Whether data messages must be received in order. If true, the messages will be sent reliably. Default true. Optional.
    Ordered *bool
    // When Ordered is false, indicates the time (in milliseconds) after which a SCTP packet will stop being retransmitted. Optional.
    MaxPacketLifeTime *int
    // When Ordered is false, indicates the maximum number of times a packet will be retransmitted. Optional.
    MaxRetransmits *int
}

// SctpParametersDump represents a dump of SCTP parameters.
type SctpParametersDump struct {
    // Port number.
    Port int
    // Initially requested number of outgoing SCTP streams.
    OS int
    // Maximum number of incoming SCTP streams.
    MIS int
    // Maximum allowed size for SCTP messages.
    MaxMessageSize int
    // Size of the send buffer.
    SendBufferSize int
    // Amount of data buffered in the SCTP socket.
    SctpBufferedAmount int
    // Whether this is a DataChannel.
    IsDataChannel bool
}
