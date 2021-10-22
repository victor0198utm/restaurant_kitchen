# Kitchen web app
Dining Hall app: <https://github.com/victor0198utm/restaurant_hall>
## Run simulation
### 1.Make network 
'''
docker network create rNet --driver bridge
'''

### 2. Build Kitchen app
From kitchen app directory: 
'''
docker build --tag kitchen .
'''

### 3. Run Kitchen container
'''
docker run -it --name kitchen --network rNet kitchen
'''

### 4. Build Hall app
From hall app directory: 
'''
docker build --tag hall .
'''

### 5. Run Hall container
'''
docker run -it --name hall --network rNet hall
'''

## Output
<img title="Applications' logs" alt="Applications' logs" src="/examples/example_2.png">
