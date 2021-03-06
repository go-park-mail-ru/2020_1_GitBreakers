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

func easyjson32ceb8acDecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit(in *jlexer.Lexer, out *RepositorySet) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		in.Skip()
		*out = nil
	} else {
		in.Delim('[')
		if *out == nil {
			if !in.IsDelim(']') {
				*out = make(RepositorySet, 0, 1)
			} else {
				*out = RepositorySet{}
			}
		} else {
			*out = (*out)[:0]
		}
		for !in.IsDelim(']') {
			var v1 Repository
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
func easyjson32ceb8acEncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit(out *jwriter.Writer, in RepositorySet) {
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
func (v RepositorySet) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson32ceb8acEncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RepositorySet) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson32ceb8acEncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RepositorySet) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson32ceb8acDecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RepositorySet) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson32ceb8acDecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit(l, v)
}
func easyjson32ceb8acDecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit1(in *jlexer.Lexer, out *Repository) {
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
		case "owner_id":
			out.OwnerID = int64(in.Int64())
		case "name":
			out.Name = string(in.String())
		case "description":
			out.Description = string(in.String())
		case "is_fork":
			out.IsFork = bool(in.Bool())
		case "created_at":
			if data := in.Raw(); in.Ok() {
				in.AddError((out.CreatedAt).UnmarshalJSON(data))
			}
		case "is_public":
			out.IsPublic = bool(in.Bool())
		case "stars":
			out.Stars = int64(in.Int64())
		case "forks":
			out.Forks = int64(in.Int64())
		case "merge_requests_open":
			out.MergeRequestsOpen = int64(in.Int64())
		case "author_login":
			out.AuthorLogin = string(in.String())
		case "parent_repository_info":
			(out.ParentRepositoryInfo).UnmarshalEasyJSON(in)
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
func easyjson32ceb8acEncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit1(out *jwriter.Writer, in Repository) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.ID))
	}
	{
		const prefix string = ",\"owner_id\":"
		out.RawString(prefix)
		out.Int64(int64(in.OwnerID))
	}
	{
		const prefix string = ",\"name\":"
		out.RawString(prefix)
		out.String(string(in.Name))
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"is_fork\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsFork))
	}
	{
		const prefix string = ",\"created_at\":"
		out.RawString(prefix)
		out.Raw((in.CreatedAt).MarshalJSON())
	}
	{
		const prefix string = ",\"is_public\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsPublic))
	}
	{
		const prefix string = ",\"stars\":"
		out.RawString(prefix)
		out.Int64(int64(in.Stars))
	}
	{
		const prefix string = ",\"forks\":"
		out.RawString(prefix)
		out.Int64(int64(in.Forks))
	}
	{
		const prefix string = ",\"merge_requests_open\":"
		out.RawString(prefix)
		out.Int64(int64(in.MergeRequestsOpen))
	}
	if in.AuthorLogin != "" {
		const prefix string = ",\"author_login\":"
		out.RawString(prefix)
		out.String(string(in.AuthorLogin))
	}
	if true {
		const prefix string = ",\"parent_repository_info\":"
		out.RawString(prefix)
		(in.ParentRepositoryInfo).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Repository) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson32ceb8acEncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Repository) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson32ceb8acEncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Repository) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson32ceb8acDecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Repository) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson32ceb8acDecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit1(l, v)
}
func easyjson32ceb8acDecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit2(in *jlexer.Lexer, out *RepoFork) {
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
		case "from_repo_id":
			out.FromRepoID = int64(in.Int64())
		case "from_author_name":
			out.FromAuthorName = string(in.String())
		case "from_repo_name":
			out.FromRepoName = string(in.String())
		case "new_name":
			out.NewName = string(in.String())
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
func easyjson32ceb8acEncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit2(out *jwriter.Writer, in RepoFork) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"from_repo_id\":"
		out.RawString(prefix[1:])
		out.Int64(int64(in.FromRepoID))
	}
	{
		const prefix string = ",\"from_author_name\":"
		out.RawString(prefix)
		out.String(string(in.FromAuthorName))
	}
	{
		const prefix string = ",\"from_repo_name\":"
		out.RawString(prefix)
		out.String(string(in.FromRepoName))
	}
	{
		const prefix string = ",\"new_name\":"
		out.RawString(prefix)
		out.String(string(in.NewName))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v RepoFork) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson32ceb8acEncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v RepoFork) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson32ceb8acEncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *RepoFork) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson32ceb8acDecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *RepoFork) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson32ceb8acDecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit2(l, v)
}
func easyjson32ceb8acDecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit3(in *jlexer.Lexer, out *ParentRepositoryInfo) {
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
			if in.IsNull() {
				in.Skip()
				out.ID = nil
			} else {
				if out.ID == nil {
					out.ID = new(int64)
				}
				*out.ID = int64(in.Int64())
			}
		case "owner_id":
			if in.IsNull() {
				in.Skip()
				out.OwnerID = nil
			} else {
				if out.OwnerID == nil {
					out.OwnerID = new(int64)
				}
				*out.OwnerID = int64(in.Int64())
			}
		case "name":
			if in.IsNull() {
				in.Skip()
				out.Name = nil
			} else {
				if out.Name == nil {
					out.Name = new(string)
				}
				*out.Name = string(in.String())
			}
		case "author_login":
			if in.IsNull() {
				in.Skip()
				out.AuthorLogin = nil
			} else {
				if out.AuthorLogin == nil {
					out.AuthorLogin = new(string)
				}
				*out.AuthorLogin = string(in.String())
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
func easyjson32ceb8acEncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit3(out *jwriter.Writer, in ParentRepositoryInfo) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ID != nil {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.Int64(int64(*in.ID))
	}
	if in.OwnerID != nil {
		const prefix string = ",\"owner_id\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(*in.OwnerID))
	}
	if in.Name != nil {
		const prefix string = ",\"name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(*in.Name))
	}
	if in.AuthorLogin != nil {
		const prefix string = ",\"author_login\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(*in.AuthorLogin))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ParentRepositoryInfo) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson32ceb8acEncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ParentRepositoryInfo) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson32ceb8acEncodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ParentRepositoryInfo) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson32ceb8acDecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ParentRepositoryInfo) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson32ceb8acDecodeGithubComGoParkMailRu20201GitBreakersInternalPkgModelsGit3(l, v)
}
