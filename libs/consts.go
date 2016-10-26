package libs

const (
	magicCookie = 0x2112A442
	fingerprint = 0x5354554e
)


const (
	errorTryAlternate                 = 300
	errorBadRequest                   = 400
	errorUnauthorized                 = 401
	errorUnassigned402                = 402
	errorForbidden                    = 403
	errorUnknownAttribute             = 420
	errorAllocationMismatch           = 437
	errorStaleNonce                   = 438
	errorUnassigned439                = 439
	errorAddressFamilyNotSupported    = 440
	errorWrongCredentials             = 441
	errorUnsupportedTransportProtocol = 442
	errorPeerAddressFamilyMismatch    = 443
	errorConnectionAlreadyExists      = 446
	errorConnectionTimeoutOrFailure   = 447
	errorAllocationQuotaReached       = 486
	errorRoleConflict                 = 487
	errorServerError                  = 500
	errorInsufficientCapacity         = 508
)
const (
	attributeFamilyIPv4 = 0x01
	attributeFamilyIPV6 = 0x02
)

const (
	AttributeMappedAddress          = 0x0001
	AttributeResponseAddress        = 0x0002
	AttributeChangeRequest          = 0x0003
	AttributeSourceAddress          = 0x0004
	AttributeChangedAddress         = 0x0005
	AttributeUsername               = 0x0006
	AttributePassword               = 0x0007
	AttributeMessageIntegrity       = 0x0008
	AttributeErrorCode              = 0x0009
	AttributeUnknownAttributes      = 0x000a
	AttributeReflectedFrom          = 0x000b
	AttributeChannelNumber          = 0x000c
	AttributeLifetime               = 0x000d
	AttributeBandwidth              = 0x0010
	AttributeXorPeerAddress         = 0x0012
	AttributeData                   = 0x0013
	AttributeRealm                  = 0x0014
	AttributeNonce                  = 0x0015
	AttributeXorRelayedAddress      = 0x0016
	AttributeRequestedAddressFamily = 0x0017
	AttributeEvenPort               = 0x0018
	AttributeRequestedTransport     = 0x0019
	AttributeDontFragment           = 0x001a
	AttributeXorMappedAddress       = 0x0020
	AttributeTimerVal               = 0x0021
	AttributeReservationToken       = 0x0022
	AttributePriority               = 0x0024
	AttributeUseCandidate           = 0x0025
	AttributePadding                = 0x0026
	AttributeResponsePort           = 0x0027
	AttributeConnectionID           = 0x002a
	AttributeXorMappedAddressExp    = 0x8020
	AttributeSoftware               = 0x8022
	AttributeAlternateServer        = 0x8023
	AttributeCacheTimeout           = 0x8027
	AttributeFingerprint            = 0x8028
	AttributeIceControlled          = 0x8029
	AttributeIceControlling         = 0x802a
	AttributeResponseOrigin         = 0x802b
	AttributeOtherAddress           = 0x802c
	AttributeEcnCheckStun           = 0x802d
	AttributeOrigin                 = 0x802f
	AttributeCiscoFlowdata          = 0xc000
	//https://tools.ietf.org/html/draft-thatcher-ice-network-cost-00
	AttributeNetworkInfo            = 0xc057

)

const (
	TypeBindingRequest                 = 0x0001
	TypeBindingResponse                = 0x0101
	TypeBindingErrorResponse           = 0x0111
	TypeSharedSecretRequest            = 0x0002
	TypeSharedSecretResponse           = 0x0102
	TypeSharedErrorResponse            = 0x0112
	TypeAllocate                       = 0x0003
	TypeAllocateResponse               = 0x0103
	TypeAllocateErrorResponse          = 0x0113
	TypeRefreshRequest                 = 0x0004
	TypeRefreshResponse                = 0x0104
	TypeRefreshErrorResponse           = 0x0114
	TypeSendIndication                 = 0x0016
	TypeDataIndication                 = 0x0017
	TypeDataResponse                   = 0x0107
	TypeDataErrorResponse              = 0x0117
	TypeCreatePermission               = 0x0008
	TypeCreatePermissionResponse       = 0x0108
	TypeCreatePermissionErrorResponse  = 0x0118
	TypeChannelBinding                 = 0x0009
	TypeChannelBindingResponse         = 0x0109
	TypeChannelBindingErrorResponse    = 0x0119
	TypeConnect                        = 0x000a
	TypeConnectResponse                = 0x010a
	TypeConnectErrorResponse           = 0x011a
	TypeConnectionBind                 = 0x000b
	TypeConnectionBindResponse         = 0x010b
	TypeConnectionBindErrorResponse    = 0x011b
	TypeConnectionAttempt              = 0x000c
	TypeConnectionAttemptResponse      = 0x010c
	TypeConnectionAttemptErrorResponse = 0x011c
)
