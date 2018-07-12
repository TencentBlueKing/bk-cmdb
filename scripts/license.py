#!/usr/bin/python 
# -*- coding: utf-8 -*-    

import sys,getopt,os

license_content='''/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except 
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and 
 * limitations under the License.
 */
 
'''

def update_license(target_file, temp_file):
    print "update: " +target_file
    with open(target_file,'r') as src_file, open(temp_file,'w') as tmp_file:
        tmp_file.write(license_content)
        is_begin=False
        for line in src_file.readlines():
            if not is_begin and not line.startswith("package"):
                continue
            is_begin=True
            tmp_file.write(line)
    os.rename(temp_file,target_file)
    os.system("gofmt -w "+temp_file +" > /dev/null 2>&1")
    

def list_dir(target_dir):
    list_dirs = os.walk(target_dir)
    for root,_,files in list_dirs:
        for f in files:
            if f.endswith(".go"):
                update_license(root+"/"+f,root+"/#"+f)

def main(argv):
    inner_dir = ''
    try:
        opts, _ = getopt.getopt(argv,"hd:",["help","dir="])
    except getopt.GetoptError:
        print 'license.py -d <directory>'
        sys.exit(2)
    if len(opts) == 0:
        print 'license.py -d <directory>'
        sys.exit(2)
    for opt, arg in opts:
        if opt in ('-h','--help'):
            print 'license.py -d <directory>'
            sys.exit()
        elif opt in ("-d", "--dir"):
            inner_dir = arg
    
    list_dir(os.path.abspath(inner_dir))

if __name__=="__main__":
    main(sys.argv[1:])
