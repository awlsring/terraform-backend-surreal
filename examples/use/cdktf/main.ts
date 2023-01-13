// Copyright (c) HashiCorp, Inc
// SPDX-License-Identifier: MPL-2.0
import { Construct } from "constructs";
import { App, HttpBackend, TerraformOutput, TerraformStack } from "cdktf";

class MyStack extends TerraformStack {
  constructor(scope: Construct, id: string) {
    super(scope, id);

    new HttpBackend(this, {
      address: "https://tf-backend.awlsring-sea.drigs.org/a/a",
      lockAddress: "https://tf-backend.awlsring-sea.drigs.org/a/a",
      unlockAddress: "https://tf-backend.awlsring-sea.drigs.org/a/a",
      username: "admin",
      password: "admin",
      skipCertVerification: false,
    })

    new TerraformOutput(this, "example", {
      value: "test",
    });
    
  }
}

const app = new App();
new MyStack(app, "cdktf");
app.synth();
