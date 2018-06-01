#!/usr/bin/python 
# -*- coding: utf-8 -*-    

import sys,getopt,os 
from string import Template

docker_file_template_str ='''
FROM  $base_image
$run_items
$copy_items
ENTRYPOINT $cwd/start.sh
'''

class DockerfileTemplate(Template):
    delimiter='$'

def generate_docker_file(image, src_dir, dst_dir, output):
    
    #print src_dir
    run_cmd_params = []
    copy_cmd_params = []
    list_dirs = os.walk(src_dir)
    for root, dirs, files in list_dirs:
        new_root = root.replace(src_dir, "")
        for d in dirs:
            target_dir = os.path.normpath(dst_dir+"/"+new_root+"/"+d)
            run_cmd_params.append("RUN mkdir -p /" + target_dir)
        for f in files:
            target_file = os.path.normpath(dst_dir+"/"+new_root+f)
            copy_cmd_params.append("COPY /"+os.path.normpath(dst_dir+"/"+new_root)+"/"+f + " /"+target_file)

    template = DockerfileTemplate(docker_file_template_str)
    result = template.substitute(dict(base_image=image, run_items="\n".join(run_cmd_params), copy_items="\n".join(copy_cmd_params),cwd="/"+dst_dir))
    
    if not os.path.exists(output):
        os.mkdir(output)    
    with open( output + "/Dockerfile."+dst_dir,'w') as tmp_file:
        tmp_file.write(result)


def list_dir(image, target, output):
    list_dirs = os.walk(target)
    for root,dirs,_ in list_dirs:
        for d in dirs:
            if d.startswith("cmdb_"):
                generate_docker_file(image, root+"/"+d, d, output)

def main(argv):
    target = ''
    image = ''
    output = ''
    try:
        opts, _ = getopt.getopt(argv,"ht:i:o:",["help","target=","base_image=","output="])
    except getopt.GetoptError:
        print 'generate.py -t <target>  -i <base_image> -o <output>'
        sys.exit(2)
    if len(opts) == 0:
        print 'generate.py -t <target> -i <base_image> -o <output>'
        sys.exit(2)
    #print opts
    for opt, arg in opts:
        if opt in ('-h','--help'):
            print 'generate.py -t <target> -i <base_image> -o <output>'
            sys.exit()
        elif opt in ("-t", "--target"):
            target = arg
        elif opt in ("-i","--base_image"):
            image = arg
        elif opt in ("-o","--output"):
            output = arg
    #print image, target
    list_dir(image, os.path.abspath(target), output)

if __name__=="__main__":
    main(sys.argv[1:])
