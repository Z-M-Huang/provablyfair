# provablyfair
Provably fair generator

Referenced from https://github.com/tyler-smith/provablyfair-dice

This is changed based on private usage.

# Example usage
```
  serverSeed, err := provablyfair.GenerateNewSeed(32)
	if err != nil {
		panic(err)
	}
	client := &provablyfair.Client{
		ServerSeed: serverSeed,
	}

  clientSeed, err := provablyfair.GenerateNewSeed(32)
	if err != nil {
		panic(err)
	}

  rollNum, serverSeed, nonce, err := client.Generate(clientSeed)
```