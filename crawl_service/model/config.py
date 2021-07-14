from typing import Optional
import os
import yaml
from pydantic import BaseModel

from crawl_service.util.root import ROOT


class VjudgeConfigModel(BaseModel):
    username: Optional[str]
    password: Optional[str]


class ConfigModel(BaseModel):
    vjudge: VjudgeConfigModel


SETTING: Optional[ConfigModel]

with open(os.path.join(ROOT, 'config.yml'), 'r', encoding='utf-8') as config_file:
    SETTING = ConfigModel(**yaml.load(config_file.read(), yaml.FullLoader))


if __name__ == '__main__':
    print(SETTING)
