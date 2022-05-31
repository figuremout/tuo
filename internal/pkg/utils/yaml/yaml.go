package yaml

import (
    "io/ioutil"
    "gopkg.in/yaml.v3"
)

func LoadYaml(path string, out interface{}) error {
    f, err := ioutil.ReadFile(path)
    if err != nil {
        return err
    }

    err = yaml.Unmarshal(f, out)
    if err != nil {
        return err
    }
    return nil
}
