from django.shortcuts import render,HttpResponse
from django.views.decorators.csrf import csrf_exempt


@csrf_exempt
def cmdb_api(request):
    if request.method == "POST":
        print(request.body)
    return HttpResponse("test")