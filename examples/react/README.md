Instead of this bad practice, I should write something like JSX for Golang to find tags and rewrite them as Golang functions.

# Sample

```
test := <Header title="test">test 
        <Link to="https://www.google.com/">Google</Link></Header>
```
Should convert to :

```
test := gox("Header", map[string]string{"title":"test"}, "test", gox("Link", map[string]string{"to":"https://www.google.com/"}) )
```