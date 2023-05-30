import { App, AppConfig } from "./app";

const port = 3000;

async function main() {
    const config: AppConfig = {
        descriptorSetPath: 'bazel-bin/example/routeguide/routeguide_proto_descriptor.pb',
        serverAddress: 'localhost:1234',
    };
    const app = new App(config);
    await app.start(port);
    const feature = await app.getFeature({ latitude: 0, longitude: 0 });
    console.log('got feature!', feature);
}

main();