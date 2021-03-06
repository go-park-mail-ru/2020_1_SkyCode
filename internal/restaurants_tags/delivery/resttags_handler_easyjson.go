// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package delivery

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson2f8b80e8DecodeGithubCom20201SkycodeInternalRestaurantsTagsDelivery(in *jlexer.Lexer, out *tagRequest) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "name":
			out.Name = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson2f8b80e8EncodeGithubCom20201SkycodeInternalRestaurantsTagsDelivery(out *jwriter.Writer, in tagRequest) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix[1:])
		out.String(string(in.Name))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v tagRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson2f8b80e8EncodeGithubCom20201SkycodeInternalRestaurantsTagsDelivery(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v tagRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson2f8b80e8EncodeGithubCom20201SkycodeInternalRestaurantsTagsDelivery(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *tagRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson2f8b80e8DecodeGithubCom20201SkycodeInternalRestaurantsTagsDelivery(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *tagRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson2f8b80e8DecodeGithubCom20201SkycodeInternalRestaurantsTagsDelivery(l, v)
}
func easyjson2f8b80e8DecodeGithubCom20201SkycodeInternalRestaurantsTagsDelivery1(in *jlexer.Lexer, out *RestTagsHandler) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson2f8b80e8EncodeGithubCom20201SkycodeInternalRestaurantsTagsDelivery1(out *jwriter.Writer, in RestTagsHandler) {
	out.RawByte('{')
	first := true
	_ = first
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RestTagsHandler) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson2f8b80e8EncodeGithubCom20201SkycodeInternalRestaurantsTagsDelivery1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RestTagsHandler) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson2f8b80e8EncodeGithubCom20201SkycodeInternalRestaurantsTagsDelivery1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RestTagsHandler) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson2f8b80e8DecodeGithubCom20201SkycodeInternalRestaurantsTagsDelivery1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RestTagsHandler) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson2f8b80e8DecodeGithubCom20201SkycodeInternalRestaurantsTagsDelivery1(l, v)
}
