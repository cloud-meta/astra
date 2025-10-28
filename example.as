resource tag {
    key: string
    value: string
}

abstract resource Resource {
    id: string
    tags: []tag
}

abstract resource VirtualMachine extents Resource {
    cpu: int
    memory: int
    storage: {
        os: FileStorage
        data: []FileStorage
    }
    network: {
        
    }
}

@provider("AWS")
@service("EC2")
@type("AWS::EC2::Instance")
resource EC2 extents VirtualMachine {
    cpu => cpu
    storage: {
        os: Volume
        data: []Volume
    }
}

abstract resource AWSResource {
    Id: string
    Trn: string
    TypeName: string
}

model CloudControlRequest<T: AWSResource> {
    TypeName: string
    ResourceModel: T
}

@auth("aws-v4")
service CloudControl<T: AWSResource> {
    region: string
    endpoint: string = "https://cloudcontrol.{region}.amazonaws.com"
    
    @post("/")
    func create(@body req: CloudControlRequest<T>)
    
    @post("/")
    func delete(@body req: CloudControlRequest<T>)
}

@cloud("aws")
provider AWS<T: AWSResource> {
    region: string
    cloudctl: CloudControl<T> = {
        region: this.region
    }

    func create(res T) {
        req: CloudControlRequest<T>{
            TypeName: meta("type")
            ResourceModel: res
        }

        cloudctl.create(req: req)
    }
}


blueprint MyArch {
    resources: {
        VirtualMachine {
            Cpu: 2
            Memory: 1024
        }
    }

    provider: AWS<?> = {
        region: us-east-1
    }

    provider: Azure<?> = {
        region: us-east-1
    }
}
