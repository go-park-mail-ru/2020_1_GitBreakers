// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package git

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

func easyjson6efd7357DecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit(in *jlexer.Lexer, out *CommitSet) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(CommitSet, 0, 1)
			} else {
				*out = CommitSet{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 Commit
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
func easyjson6efd7357EncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit(out *jwriter.Writer, in CommitSet) {
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
func (v CommitSet) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6efd7357EncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CommitSet) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6efd7357EncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CommitSet) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6efd7357DecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CommitSet) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6efd7357DecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit(l, v)
}
func easyjson6efd7357DecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit1(in *jlexer.Lexer, out *CommitRequest) {
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
		case "user_login":
			out.UserLogin = string(in.String())
		case "repo_name":
			out.RepoName = string(in.String())
		case "commit_hash":
			out.CommitHash = string(in.String())
		case "offset":
			out.Offset = int64(in.Int64())
		case "limit":
			out.Limit = int64(in.Int64())
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
func easyjson6efd7357EncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit1(out *jwriter.Writer, in CommitRequest) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"user_login\":"
		out.RawString(prefix[1:])
		out.String(string(in.UserLogin))
	}
	{
		const prefix string = ",\"repo_name\":"
		out.RawString(prefix)
		out.String(string(in.RepoName))
	}
	{
		const prefix string = ",\"commit_hash\":"
		out.RawString(prefix)
		out.String(string(in.CommitHash))
	}
	{
		const prefix string = ",\"offset\":"
		out.RawString(prefix)
		out.Int64(int64(in.Offset))
	}
	{
		const prefix string = ",\"limit\":"
		out.RawString(prefix)
		out.Int64(int64(in.Limit))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CommitRequest) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6efd7357EncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CommitRequest) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6efd7357EncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CommitRequest) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6efd7357DecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CommitRequest) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6efd7357DecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit1(l, v)
}
func easyjson6efd7357DecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit2(in *jlexer.Lexer, out *Commit) {
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
		case "commit_hash":
			out.CommitHash = string(in.String())
		case "commit_author_name":
			out.CommitAuthorName = string(in.String())
		case "commit_author_email":
			out.CommitAuthorEmail = string(in.String())
		case "commit_author_when":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.CommitAuthorWhen).UnmarshalJSON(data))
			}
		case "committer_name":
			out.CommitterName = string(in.String())
		case "committer_email":
			out.CommitterEmail = string(in.String())
		case "committer_when":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.CommitterWhen).UnmarshalJSON(data))
			}
		case "tree_hash":
			out.TreeHash = string(in.String())
		case "message":
			out.Message = string(in.String())
		case "commit_parents":
			if in.IsNull() {
				in.Skip()
				out.CommitParents = nil
			} else {
				in.Delim('[')
				if out.CommitParents == nil {
					if !in.IsDelim(']') {
						out.CommitParents = make([]string, 0, 4)
					} else {
						out.CommitParents = []string{}
					}
				} else {
					out.CommitParents = (out.CommitParents)[:0]
				}
				for !in.IsDelim(']') {
					var v4 string
					v4 = string(in.String())
					out.CommitParents = append(out.CommitParents, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
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
func easyjson6efd7357EncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit2(out *jwriter.Writer, in Commit) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"commit_hash\":"
		out.RawString(prefix[1:])
		out.String(string(in.CommitHash))
	}
	{
		const prefix string = ",\"commit_author_name\":"
		out.RawString(prefix)
		out.String(string(in.CommitAuthorName))
	}
	{
		const prefix string = ",\"commit_author_email\":"
		out.RawString(prefix)
		out.String(string(in.CommitAuthorEmail))
	}
	{
		const prefix string = ",\"commit_author_when\":"
		out.RawString(prefix)
		out.Raw((in.CommitAuthorWhen).MarshalJSON())
	}
	{
		const prefix string = ",\"committer_name\":"
		out.RawString(prefix)
		out.String(string(in.CommitterName))
	}
	{
		const prefix string = ",\"committer_email\":"
		out.RawString(prefix)
		out.String(string(in.CommitterEmail))
	}
	{
		const prefix string = ",\"committer_when\":"
		out.RawString(prefix)
		out.Raw((in.CommitterWhen).MarshalJSON())
	}
	{
		const prefix string = ",\"tree_hash\":"
		out.RawString(prefix)
		out.String(string(in.TreeHash))
	}
	{
		const prefix string = ",\"message\":"
		out.RawString(prefix)
		out.String(string(in.Message))
	}
	{
		const prefix string = ",\"commit_parents\":"
		out.RawString(prefix)
		if in.CommitParents == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.CommitParents {
				if v5 > 0 {
					out.RawByte(',')
				}
				out.String(string(v6))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Commit) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6efd7357EncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Commit) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6efd7357EncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Commit) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6efd7357DecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Commit) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6efd7357DecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit2(l, v)
}
