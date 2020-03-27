package command

import (
	"testing"

	"github.com/rentiansheng/bk_bsonmongo/writeconcern"
	"github.com/rentiansheng/bk_bsonx/network/description"
)

func TestDropDatabase(t *testing.T) {
	t.Run("Encode Write Concern for MaxWireVersion >= 5", func(t *testing.T) {
		desc := description.SelectedServer{
			Server: description.Server{
				WireVersion: &description.VersionRange{Min: 0, Max: 5},
			},
		}
		wc := writeconcern.New(writeconcern.WMajority())
		cmd := DropDatabase{WriteConcern: wc}
		write, err := cmd.encode(desc)
		noerr(t, err)
		if write.WriteConcern != wc {
			t.Error("write concern should be added to write command, but is missing")
		}
	})
	t.Run("Omit Write Concern for MaxWireVersion < 5", func(t *testing.T) {
		desc := description.SelectedServer{
			Server: description.Server{
				WireVersion: &description.VersionRange{Min: 0, Max: 4},
			},
		}
		wc := writeconcern.New(writeconcern.WMajority())
		cmd := DropDatabase{WriteConcern: wc}
		write, err := cmd.encode(desc)
		noerr(t, err)
		if write.WriteConcern != nil {
			t.Error("write concern should be omitted from write command, but is present")
		}
	})
}
