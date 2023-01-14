import { Construct } from 'constructs';
import { App, Chart, ChartProps } from 'cdk8s';
import { ConfigMap, Deployment, EnvValue, Secret, Service, ServiceType, Volume } from 'cdk8s-plus-25';
import path = require('path');

export class TerraformBackendSurrealChart extends Chart {
  constructor(scope: Construct, name: string, props?: ChartProps) {
    super(scope, name, props)
    const config = new ConfigMap(this, 'Config');
    config.addFile(path.join(__dirname, "config.yaml"));

    const configVolume = Volume.fromConfigMap(this, "config-volume", config);

    const secret = Secret.fromSecretName(this, "users", "users");
    const db = Secret.fromSecretName(this, "dbcreds", "dbcreds"); // Secret with user and password
    const usersVol = Volume.fromSecret(this, "users-volume", secret); // Secret with users

    const deployment = new Deployment(this, "deployment", {
      containers: [
        {
          image: "ghcr.io/awlsring/terraform-backend-surreal:latest",
          port: 8032,
          envVariables: {
            "CONFIG_PATH": EnvValue.fromValue("/config/config.yaml"),
            "USERS_PATH": EnvValue.fromValue("/users/users.yaml"),
            "DB_USER": EnvValue.fromSecretValue({secret: db, key: "user"}),
            "DB_PASSWORD": EnvValue.fromSecretValue({secret: db, key: "password"}),
          },
          volumeMounts: [
            {
              path: "/config",
              volume: configVolume,
            },
            {
              path: "/users",
              volume: usersVol,
            },
          ],
          securityContext: {
            allowPrivilegeEscalation: false,
            user: 1000,
          }
        }
      ],
    })

    new Service(this, "service", {
      type: ServiceType.CLUSTER_IP,
      selector: deployment,
      ports: [
        {
          port: 8032,
          targetPort: 8032,
        }
      ]
    })
  }
}

const app = new App();
new TerraformBackendSurrealChart(app, 'terraform-backend-surreal');
app.synth();
