// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package models

import (
	json "encoding/json"
	git "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
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

func easyjson315f7a6DecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModels(in *jlexer.Lexer, out *Star) {
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
		case "repo":
			out.RepoID = int64(in.Int64())
		case "vote":
			out.Vote = bool(in.Bool())
		case "author_login":
			out.AuthorLogin = string(in.String())
		case "repo_name":
			out.RepoName = string(in.String())
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
func easyjson315f7a6EncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModels(out *jwriter.Writer, in Star) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"repo\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.RepoID))
	}
	{
		const prefix string = ",\"vote\":"
		out.RawString(prefix)
		out.Bool(bool(in.Vote))
	}
	{
		const prefix string = ",\"author_login\":"
		out.RawString(prefix)
		out.String(string(in.AuthorLogin))
	}
	{
		const prefix string = ",\"repo_name\":"
		out.RawString(prefix)
		out.String(string(in.RepoName))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Star) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson315f7a6EncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModels(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Star) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson315f7a6EncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModels(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Star) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson315f7a6DecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModels(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Star) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson315f7a6DecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModels(l, v)
}
func easyjson315f7a6DecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModels1(in *jlexer.Lexer, out *RepoSet) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(RepoSet, 0, 1)
			} else {
				*out = RepoSet{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 git.Repository
			(v1).UnmarshalEasyJSON(in)
			*out = append(*out, v1)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson315f7a6EncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModels1(out *jwriter.Writer, in RepoSet) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v2, v3 := range in {
			if v2 > 0 {
				out.RawByte(',')
			}
			(v3).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v RepoSet) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson315f7a6EncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModels1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RepoSet) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson315f7a6EncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModels1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RepoSet) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson315f7a6DecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModels1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RepoSet) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson315f7a6DecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModels1(l, v)
}
func easyjson315f7a6DecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModels2(in *jlexer.Lexer, out *NewsSet) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(NewsSet, 0, 1)
			} else {
				*out = NewsSet{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v4 News
			(v4).UnmarshalEasyJSON(in)
			*out = append(*out, v4)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson315f7a6EncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModels2(out *jwriter.Writer, in NewsSet) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v5, v6 := range in {
			if v5 > 0 {
				out.RawByte(',')
			}
			(v6).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v NewsSet) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson315f7a6EncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModels2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v NewsSet) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson315f7a6EncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModels2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *NewsSet) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson315f7a6DecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModels2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *NewsSet) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson315f7a6DecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModels2(l, v)
}
func easyjson315f7a6DecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModels3(in *jlexer.Lexer, out *News) {
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
			out.ID = int64(in.Int64())
		case "author_id":
			out.AuthorID = int64(in.Int64())
		case "repo_id":
			out.RepoID = int64(in.Int64())
		case "message":
			out.Mess = string(in.String())
		case "label":
			out.Label = string(in.String())
		case "date":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.Date).UnmarshalJSON(data))
			}
		case "author_login":
			out.AuthorLogin = string(in.String())
		case "author_image":
			out.AuthorImage = string(in.String())
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
func easyjson315f7a6EncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModels3(out *jwriter.Writer, in News) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.ID))
	}
	{
		const prefix string = ",\"author_id\":"
		out.RawString(prefix)
		out.Int64(int64(in.AuthorID))
	}
	{
		const prefix string = ",\"repo_id\":"
		out.RawString(prefix)
		out.Int64(int64(in.RepoID))
	}
	{
		const prefix string = ",\"message\":"
		out.RawString(prefix)
		out.String(string(in.Mess))
	}
	{
		const prefix string = ",\"label\":"
		out.RawString(prefix)
		out.String(string(in.Label))
	}
	{
		const prefix string = ",\"date\":"
		out.RawString(prefix)
		out.Raw((in.Date).MarshalJSON())
	}
	{
		const prefix string = ",\"author_login\":"
		out.RawString(prefix)
		out.String(string(in.AuthorLogin))
	}
	{
		const prefix string = ",\"author_image\":"
		out.RawString(prefix)
		out.String(string(in.AuthorImage))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v News) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson315f7a6EncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModels3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v News) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson315f7a6EncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModels3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *News) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson315f7a6DecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModels3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *News) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson315f7a6DecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModels3(l, v)
}
func easyjson315f7a6DecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModels4(in *jlexer.Lexer, out *IssuesSet) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(IssuesSet, 0, 1)
			} else {
				*out = IssuesSet{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v7 Issue
			(v7).UnmarshalEasyJSON(in)
			*out = append(*out, v7)
			in.WantComma()
		}
		in.Delim(']')
	}
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson315f7a6EncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModels4(out *jwriter.Writer, in IssuesSet) {
	if in == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
		out.RawString("null")
	} else {
		out.RawByte('[')
		for v8, v9 := range in {
			if v8 > 0 {
				out.RawByte(',')
			}
			(v9).MarshalEasyJSON(out)
		}
		out.RawByte(']')
	}
}

// MarshalJSON supports json.Marshaler interface
func (v IssuesSet) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson315f7a6EncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModels4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v IssuesSet) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson315f7a6EncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModels4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *IssuesSet) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson315f7a6DecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModels4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *IssuesSet) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson315f7a6DecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModels4(l, v)
}
func easyjson315f7a6DecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModels5(in *jlexer.Lexer, out *Issue) {
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
			out.ID = int64(in.Int64())
		case "author_id":
			out.AuthorID = int64(in.Int64())
		case "repo_id":
			out.RepoID = int64(in.Int64())
		case "title":
			out.Title = string(in.String())
		case "message":
			out.Message = string(in.String())
		case "label":
			out.Label = string(in.String())
		case "is_closed":
			out.IsClosed = bool(in.Bool())
		case "created_at":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.CreatedAt).UnmarshalJSON(data))
			}
		case "author_login":
			out.AuthorLogin = string(in.String())
		case "author_image":
			out.AuthorImage = string(in.String())
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
func easyjson315f7a6EncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModels5(out *jwriter.Writer, in Issue) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.ID))
	}
	{
		const prefix string = ",\"author_id\":"
		out.RawString(prefix)
		out.Int64(int64(in.AuthorID))
	}
	{
		const prefix string = ",\"repo_id\":"
		out.RawString(prefix)
		out.Int64(int64(in.RepoID))
	}
	{
		const prefix string = ",\"title\":"
		out.RawString(prefix)
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"message\":"
		out.RawString(prefix)
		out.String(string(in.Message))
	}
	{
		const prefix string = ",\"label\":"
		out.RawString(prefix)
		out.String(string(in.Label))
	}
	{
		const prefix string = ",\"is_closed\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsClosed))
	}
	if true {
		const prefix string = ",\"created_at\":"
		out.RawString(prefix)
		out.Raw((in.CreatedAt).MarshalJSON())
	}
	{
		const prefix string = ",\"author_login\":"
		out.RawString(prefix)
		out.String(string(in.AuthorLogin))
	}
	{
		const prefix string = ",\"author_image\":"
		out.RawString(prefix)
		out.String(string(in.AuthorImage))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Issue) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson315f7a6EncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModels5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Issue) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson315f7a6EncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModels5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Issue) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson315f7a6DecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModels5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Issue) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson315f7a6DecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModels5(l, v)
}
