package match

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_insertNode(t *testing.T) {
	assert := assert.New(t)

	// insert first node
	newNode := &node{content: &SinglePerson{Name: "James", Height: 180, Gender: male, WantedDate: 1}}
	root := insertNode(nil, newNode)
	assert.Equal(newNode, root)

	// insert to right
	newNode = &node{content: &SinglePerson{Name: "John", Height: 190, Gender: male, WantedDate: 2}}
	root = insertNode(root, newNode)
	assert.Equal(newNode, root.right)

	// insert to left
	newNode = &node{content: &SinglePerson{Name: "Jane", Height: 165, Gender: female, WantedDate: 3}}
	root = insertNode(root, newNode)
	assert.Equal(newNode, root.Left)
}

func Test_findMatchesWithLimitAndCondition(t *testing.T) {
	simpeNodes := []*node{
		{content: &SinglePerson{Name: "James", Height: 180, Gender: male, WantedDate: 1}},
		{content: &SinglePerson{Name: "John", Height: 190, Gender: male, WantedDate: 2}},
		{content: &SinglePerson{Name: "Jane", Height: 165, Gender: female, WantedDate: 3}},
	}

	var simpleRoot *node
	for _, n := range simpeNodes {
		simpleRoot = insertNode(simpleRoot, n)
	}

	type args struct {
		root  *node
		match *node
		limit int
	}
	tests := []struct {
		name  string
		args  args
		want  []*node
		want1 int
	}{
		// TODO: Add test cases.
		{
			name: "General",
			args: args{
				root:  simpleRoot,
				match: simpeNodes[0],
				limit: 1,
			},
			want: []*node{
				simpeNodes[2],
			},
			want1: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := findMatchesWithLimitAndCondition(tt.args.root, tt.args.match, tt.args.limit)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findMatchesWithLimitAndCondition() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("findMatchesWithLimitAndCondition() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
