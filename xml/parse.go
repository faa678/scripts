package localxml

import (
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"scripts/collections"
)

var (
	ErrorInvalidXmlSyntax = errors.New("invalid xml syntax")
)

var (
	// only support some simple xml style like `<a>aaa</a>`
	spaLineReg = regexp.MustCompile(`^(\s*)$`)
	headTagReg = regexp.MustCompile(`^<\w+?>`)
	contentReg = regexp.MustCompile(`^[\w:/\.,-]+`)
	endTagReg  = regexp.MustCompile(`^</\w+?>`)
)

type Node struct {
	Name string
	Val  string
	Par  *Node
	Chi  map[string]*Node
}

// Parse parses the docoder into a prefix-tree
func Parse(decoder *Decoder) (*Node, error) {
	var root, parNode, curNode *Node = nil, nil, nil
	stack := collections.NewStack(0)

	for _, token := range decoder.Tokens {
		val := token.Value
		switch token.Type {
		case TokenTypeHTag:
			node := &Node{
				Name: val,
				Par:  parNode,
			}
			curNode = node
			if !stack.IsEmpty() {
				parNode = stack.Peek().(*Node)
			} else {
				root = curNode
			}
			if parNode != nil {
				if parNode.Chi == nil {
					parNode.Chi = make(map[string]*Node)
				}
				if _, ok := parNode.Chi[val]; ok {
					return nil, ErrorInvalidXmlSyntax
				}
				parNode.Chi[val] = node
			}
			curNode.Par = parNode
			_ = stack.Offer(curNode)
		case TokenTypeContent:
			val := string(val)
			curNode.Val = val
		case TokenTypeETag:
			node, err := stack.Poll()
			if err != nil {
				return nil, ErrorInvalidXmlSyntax
			}
			if node == parNode {
				parNode = parNode.Par
			}
		default:
			return nil, errors.New("unknown token")
		}

	}

	return root, nil
}

type Token struct {
	Type  TokenType
	Value string
}

type Decoder struct {
	Tokens []*Token
	Size   int32
}

// Decode decodes the xml contents and validate the syntax
func Decode(filename string) (*Decoder, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	contents := string(bytes)
	ls := strings.Split(contents, "\n")

	tokens := make([]*Token, 0)
	stack := new(collections.Stack)
	var curToken *Token
	for _, line := range ls {
		if spaLineReg.MatchString(line) {
			continue
		}
		line = strings.TrimSpace(line)
		// fmt.Println(line)
		l := len(line)
		for idx := 0; idx < l; {
			bytes := []byte(line[idx:])
			peek := &Token{}
			if !stack.IsEmpty() {
				peek = stack.Peek().(*Token)
			}
			switch t, i, j := matchTokenType(bytes); t {
			case TokenTypeHTag:
				if !stack.IsEmpty() && peek.Type == TokenTypeContent {
					// fmt.Println("h tag")
					return nil, ErrorInvalidXmlSyntax
				}
				curToken = &Token{
					Type:  TokenTypeHTag,
					Value: string(bytes[i+1 : j-1]),
				}
				_ = stack.Offer(curToken)
				idx += j
			case TokenTypeContent:
				if stack.IsEmpty() || peek.Type == TokenTypeContent || peek.Type == TokenTypeETag {
					// fmt.Println("content")
					return nil, ErrorInvalidXmlSyntax
				}
				curToken = &Token{
					Type:  TokenTypeContent,
					Value: string(bytes[i:j]),
				}
				_ = stack.Offer(curToken)
				idx += j
			case TokenTypeETag:
				if stack.IsEmpty() {
					// fmt.Println("e tag")
					return nil, ErrorInvalidXmlSyntax
				}
				curToken = &Token{
					Type:  TokenTypeETag,
					Value: string(bytes[i+2 : j-1]),
				}
				p, err := stack.Poll()
				if err != nil {
					return nil, ErrorInvalidXmlSyntax
				}
				token := p.(*Token)
				if token.Type == TokenTypeContent {
					p, _ = stack.Poll()
					token = p.(*Token)
				}
				// fmt.Println(token.Value, curToken.Value)
				if token.Value != curToken.Value {
					return nil, ErrorInvalidXmlSyntax
				}
				idx += j
			default:
				return nil, ErrorInvalidXmlSyntax
			}
			// fmt.Println(curToken.Value)
			tokens = append(tokens, curToken)
		}
	}
	if !stack.IsEmpty() {
		return nil, ErrorInvalidXmlSyntax
	}
	decoder := &Decoder{
		Tokens: tokens,
		Size:   int32(len(tokens)),
	}
	return decoder, nil
}

type TokenType int

const (
	TokenTypeUnknown = 0
	TokenTypeHTag    = 1
	TokenTypeContent = 2
	TokenTypeETag    = 3
	TokenTypeComment = 4
)

func matchTokenType(bytes []byte) (TokenType, int, int) {
	if idxs := headTagReg.FindIndex(bytes); len(idxs) == 2 {
		// fmt.Println(string(bytes), " aaa ", idxs[0], idxs[1])
		return TokenTypeHTag, idxs[0], idxs[1]
	}
	if idxs := contentReg.FindIndex(bytes); len(idxs) == 2 {
		// fmt.Println(string(bytes)+" bbb", idxs[0], idxs[1])
		return TokenTypeContent, idxs[0], idxs[1]
	}
	if idxs := endTagReg.FindIndex(bytes); len(idxs) == 2 {
		// fmt.Println(string(bytes)+" ccc", idxs[0], idxs[1])
		return TokenTypeETag, idxs[0], idxs[1]
	}
	fmt.Println("unknow tag")
	return TokenTypeUnknown, 0, 0
}
