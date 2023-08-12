package migrations

import "embed"

// setting the global variable that has the file system that we need
// so its going to have the embedding that we need
//  we need the set the comment above the variable to set
// which files that going to be embedded

//go:embed *.sql
var FS embed.FS
