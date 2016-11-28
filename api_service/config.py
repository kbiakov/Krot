'''
This script is necessary if you want to create config file automatically.

Run script with following params:
-a : will set admin level permissions to specified user (even if not exist)
-e : sent email, logs & usage statistics on this email (shoul be verified later)

Author: Kirill Biakov (kbiakov@gmail.com)
'''

import hashlib, random

class Config:
    def __init__(self):
        self.admin = 'admin'
        self.email = 'admin@gmail.com'
        self.api_key = self.generate_api_key()

    def generate_api_key(self):
        rnd_str = str(random.getrandbits(256))
        return hashlib.sha224(rnd_str).hexdigest()
    
    def to_JSON(self):
        return json.dumps(self, default = lambda o: o.__dict__,
            sort_keys = True, ensure_ascii = False, indent = 4)

# at this point read config by presented args

import sys, getopt

def config_with_args(argv):
    config = Config()

    try:
        opts, args = getopt.getopt(argv, 'ha:e:', ['admin=', 'email='])
    except getopt.GetoptError:
        print 'init_config.py -a <admin> -e <email>'
        sys.exit(2)

    for opt, arg in opts:
        if opt == '-h':
            print 'init_config.py -a <admin> -e <email>'
            sys.exit()
        elif opt in ('-a', '--admin'):
            config.admin = arg
        elif opt in ('-e', '--email'):
            config.email = arg
    
    return config.to_JSON()

# finally, create config file

import codecs, json

with codecs.open('config.json', 'w', 'utf8') as f:
     f.write(config_with_args(sys.argv[1:]))
