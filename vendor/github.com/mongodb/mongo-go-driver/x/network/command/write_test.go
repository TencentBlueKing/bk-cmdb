package command

import (
	"bytes"
	"testing"

	"github.com/mongodb/mongo-go-driver/mongo/writeconcern"
	"github.com/mongodb/mongo-go-driver/x/bsonx"
	"github.com/mongodb/mongo-go-driver/x/bsonx/bsoncore"
	"github.com/mongodb/mongo-go-driver/x/network/description"
	"github.com/mongodb/mongo-go-driver/x/network/wiremessage"
)

func TestWrite(t *testing.T) {
	t.Run("Encode", func(t *testing.T) {
		t.Run("should not encode empty write concern", func(t *testing.T) {
			cmd := bsonx.Doc{{"fakeCommand", bsonx.Int32(1)}}
			want, err := append(cmd, bsonx.Elem{"$db", bsonx.String("foobar")}).MarshalBSON()
			noerr(t, err)
			w := Write{
				DB:           "foobar",
				Command:      cmd,
				WriteConcern: writeconcern.New(),
			}
			wm, err := w.Encode(description.SelectedServer{
				Server: description.Server{
					WireVersion: &description.VersionRange{Min: 0, Max: wiremessage.OpmsgWireVersion},
				},
			})
			noerr(t, err)
			msg, ok := wm.(wiremessage.Msg)
			if !ok {
				t.Errorf("Expected an OP_MSG wire message, but got something else. got %v", wm)
			}
			got := msg.Sections[0].(wiremessage.SectionBody).Document
			if !bytes.Equal(got, want) {
				t.Errorf("Command documents do not match. got %v; want %v", bsoncore.Document(got), bsoncore.Document(want))
			}
		})
	})
}
