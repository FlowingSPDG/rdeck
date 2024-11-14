package connector

import "context"

// connector connects Input and Output, and makes it easier to handle

type Connector interface {
	Start(ctx context.Context) error
}
