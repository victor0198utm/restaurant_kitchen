# Kitchen web app


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
docker run -dit --name kitchen --network rNet kitchen
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
After starting the Hall container, it will make 10 requests to the kitchen. This is logged in container's std output. The Kitchen container also logs the receiving of the requests and the made requests.

<img title="Applications' logs" alt="Applications' logs" src="/examples/example_1.png">
