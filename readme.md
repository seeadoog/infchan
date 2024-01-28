
## InfChan
InfChan is a go chan with infinity capacity. Put elem to  Chan will never block.


### example
````

c := NewChan[int](1)

for i := 0; i < 10; i++ {
    c.Put(i)
}
c.Close()
for {

    data, ok := <-c.Get()
    fmt.Println(data, ok)
    if !ok {
     break
    }
}

````
