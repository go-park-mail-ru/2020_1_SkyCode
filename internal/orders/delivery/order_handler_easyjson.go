// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package delivery

import (
	json "encoding/json"
	models "github.com/2020_1_Skycode/internal/models"
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

func easyjson83df9c93DecodeGithubCom20201SkycodeInternalOrdersDelivery(in *jlexer.Lexer, out *orderRequest) {
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
		case "userId":
			out.UserID = uint64(in.Uint64())
		case "restId":
			out.RestID = uint64(in.Uint64())
		case "address":
			out.Address = string(in.String())
		case "comment":
			out.Comment = string(in.String())
		case "phone":
			out.Phone = string(in.String())
		case "personNum":
			out.PersonNum = uint32(in.Uint32())
		case "products":
			if in.IsNull() {
				in.Skip()
				out.Products = nil
			} else {
				in.Delim('[')
				if out.Products == nil {
					if !in.IsDelim(']') {
						out.Products = make([]*models.OrderProduct, 0, 8)
					} else {
						out.Products = []*models.OrderProduct{}
					}
				} else {
					out.Products = (out.Products)[:0]
				}
				for !in.IsDelim(']') {
					var v1 *models.OrderProduct
					if in.IsNull() {
						in.Skip()
						v1 = nil
					} else {
						if v1 == nil {
							v1 = new(models.OrderProduct)
						}
						easyjson83df9c93DecodeGithubCom20201SkycodeInternalModels(in, v1)
					}
					out.Products = append(out.Products, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "price":
			out.Price = float32(in.Float32())
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

// UnmarshalJSON supports json.Unmarshaler interface
func (v *orderRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson83df9c93DecodeGithubCom20201SkycodeInternalOrdersDelivery(&r, v)
	return r.Error()
}

func easyjson83df9c93DecodeGithubCom20201SkycodeInternalModels(in *jlexer.Lexer, out *models.OrderProduct) {
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
		case "id":
			out.ID = uint64(in.Uint64())
		case "orderId":
			out.OrderID = uint64(in.Uint64())
		case "productId":
			out.ProductID = uint64(in.Uint64())
		case "count":
			out.Count = uint32(in.Uint32())
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
func easyjson83df9c93EncodeGithubCom20201SkycodeInternalModels(out *jwriter.Writer, in models.OrderProduct) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Uint64(uint64(in.ID))
	}
	{
		const prefix string = ",\"orderId\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.OrderID))
	}
	{
		const prefix string = ",\"productId\":"
		out.RawString(prefix)
		out.Uint64(uint64(in.ProductID))
	}
	{
		const prefix string = ",\"count\":"
		out.RawString(prefix)
		out.Uint32(uint32(in.Count))
	}
	out.RawByte('}')
}
