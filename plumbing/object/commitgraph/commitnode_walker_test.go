package commitgraph

import (
	"strings"
	"testing"

	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/cache"
	commitgraph "github.com/go-git/go-git/v6/plumbing/format/commitgraph"
	"github.com/go-git/go-git/v6/plumbing/format/packfile"
	"github.com/go-git/go-git/v6/storage/filesystem"
	"github.com/stretchr/testify/assert"

	fixtures "github.com/go-git/go-git-fixtures/v5"
)

func TestCommitNodeIter(t *testing.T) {
	t.Parallel()

	f := fixtures.ByTag("commit-graph-chain-2").One()

	storer := newUnpackRepository(f)

	index, err := commitgraph.OpenChainOrFileIndex(storer.Filesystem())
	assert.NoError(t, err)

	nodeIndex := NewGraphCommitNodeIndex(index, storer)

	head, err := nodeIndex.Get(plumbing.NewHash("ec6f456c0e8c7058a29611429965aa05c190b54b"))
	assert.NoError(t, err)

	testTopoOrder(t, head)
	testDateOrder(t, head)
	testAuthorDateOrder(t, head)
}

func newUnpackRepository(f *fixtures.Fixture) *filesystem.Storage {
	storer := filesystem.NewStorage(f.DotGit(), cache.NewObjectLRUDefault())
	p := f.Packfile()
	defer p.Close()
	packfile.UpdateObjectStorage(storer, p)
	return storer
}

func testTopoOrder(t *testing.T, head CommitNode) {
	iter := NewCommitNodeIterTopoOrder(
		head,
		nil,
		nil,
	)

	var commits []string
	iter.ForEach(func(c CommitNode) error {
		commits = append(commits, c.ID().String())
		return nil
	})

	assert.Equal(t, commits, strings.Split(`ec6f456c0e8c7058a29611429965aa05c190b54b
d82f291cde9987322c8a0c81a325e1ba6159684c
3048d280d2d5b258d9e582a226ff4bbed34fd5c9
27aa8cdd2431068606741a589383c02c149ea625
fa058d42fa3bc53f39108a56dad67157169b2191
6c629843a1750a27c9af01ed2985f362f619c47a
d10a0e7c1f340a6cfc14540a5f8c508ce7e2eabf
d0a18ccd8eea3bdabc76d6dc5420af1ea30aae9f
cf2874632223220e0445abf0a7806dc772c0b37a
758ac33217f092bfcded4ad4774954ac054c9609
214e1dca024fb6da5ed65564d2de734df5dc2127
70923099e61fa33f0bc5256d2f938fa44c4df10e
bcaa1ac5644b16f1febb72f31e204720b7bb8934
e1d8866ffa78fa16d2f39b0ba5344a7269ee5371
2275fa7d0c75d20103f90b0e1616937d5a9fc5e6
bdd9a92789d4a86b20a8d3df462df373f41acf23
b359f11ea09e642695edcd114b463da4395b10c1
6f43e8933ba3c04072d5d104acc6118aac3e52ee
ccafe8bd5f9dbfb8b98b0da03ced29608dcfdeec
939814f341fdd5d35e81a3845a33c4fedb19d2d2
5f5ad88bf2babe506f927d64d2b7a1e1493dc2ae
a2014124ca3b3f9ff28fbab0a83ce3c71bf4622e
77906b653c3eb8a1cd5bd7254e161c00c6086d83
465cba710284204f9851854587c2887c247222db
b9471b13256703d3f5eb88b280b4a16ce325ec1b
62925030859646daeeaf5a4d386a0c41e00dda8a
5f56aea0ca8b74215a5b982bca32236e1e28c76b
23148841baa5dbce48f6adcb7ddf83dcd97debb3
c336d16298a017486c4164c40f8acb28afe64e84
31eae7b619d166c366bf5df4991f04ba8cebea0a
d2a38b4a5965d529566566640519d03d2bd10f6c
b977a025ca21e3b5ca123d8093bd7917694f6da7
35b585759cbf29f8ec428ef89da20705d59f99ec
c2bbf9fe8009b22d0f390f3c8c3f13937067590f
fc9f0643b21cfe571046e27e0c4565f3a1ee96c8
c088fd6a7e1a38e9d5a9815265cb575bb08d08ff
5fddbeb678bd2c36c5e5c891ab8f2b143ced5baf
5d7303c49ac984a9fec60523f2d5297682e16646`, "\n"))
}

func testDateOrder(t *testing.T, head CommitNode) {
	iter := NewCommitNodeIterDateOrder(
		head,
		nil,
		nil,
	)

	var commits []string
	iter.ForEach(func(c CommitNode) error {
		commits = append(commits, c.ID().String())
		return nil
	})

	assert.Equal(t, commits, strings.Split(`ec6f456c0e8c7058a29611429965aa05c190b54b
3048d280d2d5b258d9e582a226ff4bbed34fd5c9
d82f291cde9987322c8a0c81a325e1ba6159684c
27aa8cdd2431068606741a589383c02c149ea625
fa058d42fa3bc53f39108a56dad67157169b2191
d0a18ccd8eea3bdabc76d6dc5420af1ea30aae9f
6c629843a1750a27c9af01ed2985f362f619c47a
cf2874632223220e0445abf0a7806dc772c0b37a
d10a0e7c1f340a6cfc14540a5f8c508ce7e2eabf
758ac33217f092bfcded4ad4774954ac054c9609
214e1dca024fb6da5ed65564d2de734df5dc2127
70923099e61fa33f0bc5256d2f938fa44c4df10e
bcaa1ac5644b16f1febb72f31e204720b7bb8934
e1d8866ffa78fa16d2f39b0ba5344a7269ee5371
2275fa7d0c75d20103f90b0e1616937d5a9fc5e6
bdd9a92789d4a86b20a8d3df462df373f41acf23
b359f11ea09e642695edcd114b463da4395b10c1
6f43e8933ba3c04072d5d104acc6118aac3e52ee
ccafe8bd5f9dbfb8b98b0da03ced29608dcfdeec
939814f341fdd5d35e81a3845a33c4fedb19d2d2
5f5ad88bf2babe506f927d64d2b7a1e1493dc2ae
a2014124ca3b3f9ff28fbab0a83ce3c71bf4622e
77906b653c3eb8a1cd5bd7254e161c00c6086d83
465cba710284204f9851854587c2887c247222db
b9471b13256703d3f5eb88b280b4a16ce325ec1b
62925030859646daeeaf5a4d386a0c41e00dda8a
5f56aea0ca8b74215a5b982bca32236e1e28c76b
23148841baa5dbce48f6adcb7ddf83dcd97debb3
c336d16298a017486c4164c40f8acb28afe64e84
31eae7b619d166c366bf5df4991f04ba8cebea0a
b977a025ca21e3b5ca123d8093bd7917694f6da7
d2a38b4a5965d529566566640519d03d2bd10f6c
35b585759cbf29f8ec428ef89da20705d59f99ec
c2bbf9fe8009b22d0f390f3c8c3f13937067590f
fc9f0643b21cfe571046e27e0c4565f3a1ee96c8
c088fd6a7e1a38e9d5a9815265cb575bb08d08ff
5fddbeb678bd2c36c5e5c891ab8f2b143ced5baf
5d7303c49ac984a9fec60523f2d5297682e16646`, "\n"))
}

func testAuthorDateOrder(t *testing.T, head CommitNode) {
	iter := NewCommitNodeIterAuthorDateOrder(
		head,
		nil,
		nil,
	)

	var commits []string
	iter.ForEach(func(c CommitNode) error {
		commits = append(commits, c.ID().String())
		return nil
	})

	assert.Equal(t, commits, strings.Split(`ec6f456c0e8c7058a29611429965aa05c190b54b
3048d280d2d5b258d9e582a226ff4bbed34fd5c9
d82f291cde9987322c8a0c81a325e1ba6159684c
27aa8cdd2431068606741a589383c02c149ea625
fa058d42fa3bc53f39108a56dad67157169b2191
d0a18ccd8eea3bdabc76d6dc5420af1ea30aae9f
6c629843a1750a27c9af01ed2985f362f619c47a
cf2874632223220e0445abf0a7806dc772c0b37a
d10a0e7c1f340a6cfc14540a5f8c508ce7e2eabf
758ac33217f092bfcded4ad4774954ac054c9609
214e1dca024fb6da5ed65564d2de734df5dc2127
70923099e61fa33f0bc5256d2f938fa44c4df10e
bcaa1ac5644b16f1febb72f31e204720b7bb8934
e1d8866ffa78fa16d2f39b0ba5344a7269ee5371
2275fa7d0c75d20103f90b0e1616937d5a9fc5e6
bdd9a92789d4a86b20a8d3df462df373f41acf23
b359f11ea09e642695edcd114b463da4395b10c1
6f43e8933ba3c04072d5d104acc6118aac3e52ee
ccafe8bd5f9dbfb8b98b0da03ced29608dcfdeec
939814f341fdd5d35e81a3845a33c4fedb19d2d2
5f5ad88bf2babe506f927d64d2b7a1e1493dc2ae
a2014124ca3b3f9ff28fbab0a83ce3c71bf4622e
77906b653c3eb8a1cd5bd7254e161c00c6086d83
465cba710284204f9851854587c2887c247222db
b9471b13256703d3f5eb88b280b4a16ce325ec1b
5f56aea0ca8b74215a5b982bca32236e1e28c76b
62925030859646daeeaf5a4d386a0c41e00dda8a
23148841baa5dbce48f6adcb7ddf83dcd97debb3
c336d16298a017486c4164c40f8acb28afe64e84
31eae7b619d166c366bf5df4991f04ba8cebea0a
b977a025ca21e3b5ca123d8093bd7917694f6da7
d2a38b4a5965d529566566640519d03d2bd10f6c
35b585759cbf29f8ec428ef89da20705d59f99ec
c2bbf9fe8009b22d0f390f3c8c3f13937067590f
fc9f0643b21cfe571046e27e0c4565f3a1ee96c8
c088fd6a7e1a38e9d5a9815265cb575bb08d08ff
5fddbeb678bd2c36c5e5c891ab8f2b143ced5baf
5d7303c49ac984a9fec60523f2d5297682e16646`, "\n"))
}
