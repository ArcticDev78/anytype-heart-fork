package block

import (
	"github.com/anytypeio/any-sync/commonspace/object/accountdata"
	"github.com/anytypeio/any-sync/commonspace/object/acl/list"
	"github.com/anytypeio/any-sync/commonspace/object/tree/objecttree"
	"github.com/anytypeio/any-sync/commonspace/object/tree/treechangeproto"
	"github.com/anytypeio/any-sync/util/crypto"
	spaceservice "github.com/anytypeio/go-anytype-middleware/space"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_Payloads(t *testing.T) {
	// doing some any-sync preparations
	changePayload := []byte("some")
	keys, err := accountdata.NewRandom()
	require.NoError(t, err)
	aclList, err := list.NewTestDerivedAcl("spaceId", keys)
	require.NoError(t, err)
	timestamp := time.Now().Add(time.Hour).Unix()

	checkRoot := func(root *treechangeproto.RawTreeChangeWithId, changePayload []byte, changeType string, timestamp int64) {
		builder := objecttree.NewChangeBuilder(crypto.NewKeyStorage(), root)
		ch, err := builder.Unmarshall(root, true)
		require.NoError(t, err)
		rootModel := &treechangeproto.TreeChangeInfo{}
		err = proto.Unmarshal(ch.Data, rootModel)
		require.NoError(t, err)

		require.Equal(t, rootModel.ChangePayload, changePayload)
		require.Equal(t, rootModel.ChangeType, spaceservice.ChangeType)
		require.Equal(t, ch.Timestamp, timestamp)
	}

	t.Run("test create payload", func(t *testing.T) {
		firstPayload, err := createPayload("spaceId", keys.SignKey, changePayload, timestamp)
		require.NoError(t, err)
		firstRoot, err := objecttree.CreateObjectTreeRoot(firstPayload, aclList)
		require.NoError(t, err)

		secondPayload, err := createPayload("spaceId", keys.SignKey, changePayload, timestamp)
		require.NoError(t, err)
		secondRoot, err := objecttree.CreateObjectTreeRoot(secondPayload, aclList)
		require.NoError(t, err)

		// checking that created roots are not equal
		require.NotEqual(t, firstRoot, secondRoot)

		checkRoot(firstRoot, changePayload, spaceservice.ChangeType, timestamp)
		checkRoot(secondRoot, changePayload, spaceservice.ChangeType, timestamp)
	})

	t.Run("test derive payload", func(t *testing.T) {
		firstPayload := derivePayload("spaceId", keys.SignKey, changePayload)
		firstRoot, err := objecttree.CreateObjectTreeRoot(firstPayload, aclList)
		require.NoError(t, err)

		secondPayload := derivePayload("spaceId", keys.SignKey, changePayload)
		secondRoot, err := objecttree.CreateObjectTreeRoot(secondPayload, aclList)
		require.NoError(t, err)

		// checking that derived roots are equal
		require.Equal(t, firstRoot, secondRoot)
		checkRoot(firstRoot, changePayload, spaceservice.ChangeType, 0)
	})
}