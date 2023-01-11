// Copyright (c) HashiCorp, Inc
// SPDX-License-Identifier: MPL-2.0
import { Construct } from "constructs";
import { App, HttpBackend, TerraformOutput, TerraformStack } from "cdktf";

class MyStack extends TerraformStack {
  constructor(scope: Construct, id: string) {
    super(scope, id);

    new HttpBackend(this, {
      address: "http://localhost:8032/a/a",
      lockAddress: "http://localhost:8032/a/a",
      unlockAddress: "http://localhost:8032/a/a",
      username: "admin",
      password: "admin",
      skipCertVerification: true,
    })

    new TerraformOutput(this, "example", {
      value: "test",
    });
    
  }
}

const app = new App();
new MyStack(app, "cdktf");
app.synth();
