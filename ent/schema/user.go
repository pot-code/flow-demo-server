package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("username").Comment("登录名").Unique(),
		field.String("name").Comment("姓名"),
		field.String("mobile").Comment("手机号").Unique(),
		field.String("password").Comment("密码").Unique(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
