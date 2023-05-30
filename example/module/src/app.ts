
import express, { Application, Request, Response } from 'express';
import { PackageDefinition, Options, loadFileDescriptorSetFromBuffer } from '@grpc/proto-loader';
import { Message, Root, RootConstructor } from 'protobufjs';
import { promises as fsPromises } from 'fs';
import { IFileDescriptorSet } from 'protobufjs/ext/descriptor';
import * as grpc from '@grpc/grpc-js'

type DecodedDescriptorSet = Message<IFileDescriptorSet> & IFileDescriptorSet;

type Point = {
    latitude: number;
    longitude: number;
};

type Feature = {
    name: string;
    location: Point;
};

export type AppConfig = {
    descriptorSetPath: string;
    serverAddress: string;
}

export class App {
    private express: Application;
    private client: grpc.Client;
    private grpc: grpc.GrpcObject;
    // private root: Root;
    private protoPkg: PackageDefinition;

    constructor(private config: AppConfig) {
        this.express = express();
    }

    handleGetRoot(req: Request, res: Response) {
        res.send('Hello World!');
    }

    async start(port: number) {

        this.express.get('/', this.handleGetRoot.bind(this));

        const creds: grpc.ChannelCredentials = grpc.credentials.createInsecure();
        const clientOptions: grpc.ClientOptions = {};
        this.client = new grpc.Client(this.config.serverAddress, creds, {

        });

        this.grpc = await this.loadGrpcObject();
        // console.log('grpc', this.grpc);
        // console.log('root', this.root);
        // @ts-ignore
        // console.log('point', this.grpc.example.routeguide.Point);

        this.express.listen(port, () => {
            return console.log(`App is listening at http://localhost:${port}`);
        });

    }

    async getFeature(point: Point): Promise<Feature> {
        const method = this.protoPkg['example.routeguide.RouteGuide']['GetFeature'] as grpc.MethodDefinition<any, any>;
        // console.log('method', method);
        const md: grpc.Metadata = new grpc.Metadata();
        const options: grpc.CallOptions = {};

        return new Promise((resolve, reject) => {
            const callback = (err: grpc.ServiceError, value: Feature) => {
                console.log('getFeature callback:', err, value);
                if (err) {
                    reject(err);
                } else {
                    resolve(value);
                }
            };
            const call: grpc.ClientUnaryCall = this.client.makeUnaryRequest(
                method.path,
                method.requestSerialize,
                method.responseDeserialize,
                point,
                md,
                options,
                callback,
            );
            call.on('status', (md: grpc.Metadata) => {
                console.log('getFeature status callback:', md);
            });
            call.on('metadata', (md: grpc.Metadata) => {
                console.log('getFeature metadata callback:', md);
            });
        });
    }

    async loadGrpcObject(): Promise<grpc.GrpcObject> {
        const descriptor = require('protobufjs/ext/descriptor');
        const descriptorSetBuffer: Buffer = await fsPromises.readFile(this.config.descriptorSetPath);
        const decodedDescriptorSet = descriptor.FileDescriptorSet.decode(descriptorSetBuffer) as DecodedDescriptorSet;
        const options: Options = {};

        const protoPkg: PackageDefinition = loadFileDescriptorSetFromBuffer(descriptorSetBuffer, options);

        // const root: Root = new Root();
        // const descriptorSetRoot = (Root as RootConstructor).fromDescriptor(decodedDescriptorSet);
        // root.add(descriptorSetRoot);
        // this.root = root;
        this.protoPkg = protoPkg;

        return grpc.loadPackageDefinition(protoPkg);
    }
}
