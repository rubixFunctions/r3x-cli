// Copyright © 2019 RubiXFunctions
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
)

func InitializePyFunction(function *Function, schema *Schema){
	if !exists(function.AbsPath()) {
		err := os.MkdirAll(function.AbsPath(), os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else if !isEmpty(function.AbsPath()) {
		fmt.Println("Function can not be bootstrapped in a non empty direcctory: " + function.AbsPath())
		return
	}

	createPyDockerfile(function)
	createPyMain(function)
	createSchema(schema, function)
	fmt.Println(`Your Function is ready at` + function.AbsPath())
}

func createPyDockerfile(function *Function){
	dockerTemplate := `FROM python:2.7.15-alpine

ADD r3x-func.py /

RUN python -m pip install -i https://test.pypi.org/simple/ r3x 

ENV PORT=8080

CMD [ "python", "./r3x-func.py" ]`

	createFile(function, dockerTemplate, "Dockerfile")
}

func createPyMain(function *Function){
	pyTemplate := `import r3x
import json

def r3xFunc(input):
    i = json.loads(input)
    for key,value in i.items():
        if str(key) == "name":
            res = {"message": "hello {}".format(value)}
        else:
            res = {"message" : "hello r3x"}
    json_res = json.dumps(res)
    return json_res

if __name__ == "__main__":  
    r3x.execute(r3xFunc)`

	createFile(function, pyTemplate, "r3x-func.py")
}


