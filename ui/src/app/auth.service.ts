import { Injectable } from '@angular/core';
import { Router } from '@angular/router';

import * as auth0 from 'auth0-js';
import { environment } from 'src/environments/environment';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  constructor(public router: Router) { }

  accessToken: string;
  idToken: string;
  expiresAt: string;

  auth0 = new auth0.WebAuth({
    clientID: environment.clientId,
    domain: environment.domain,
    responseType: 'token id_token',
    audience: environment.audience,
    redirectUri: environment.callback,
    scope: 'openid profile email'
  });

  public login(): void {
    this.auth0.authorize();
  }

  public handleAuthentication(): void {
    this.auth0.parseHash((err, authResult) => {
      if (err) { console.log(err); }
      if (!err && authResult && authResult.accessToken && authResult.idToken) {
        window.location.hash = '';
        this.setSession(authResult);
      }

      this.router.navigate(['/home']);

    });
  }

  private setSession(authResult: auth0.Auth0DecodedHash) {
    // Set time access token will expire
    this.expiresAt = JSON.stringify((authResult.expiresIn * 1000) + new Date().getTime());
    this.accessToken = authResult.accessToken;
    this.idToken = authResult.idToken;
  }

  public logout(): void {
    this.accessToken = null;
    this.idToken = null;
    this.expiresAt = null;
    this.router.navigate(['/']);
    // this.auth0.logout({
    //   returnTo: "http://localhost:4200/",
    //   clientID: environment.clientId
    // });
  }

  public isAuthenticated(): boolean {
    // Check if current time is past access token expiration time
    const expiresAt = JSON.parse(this.expiresAt || '{}');
    return new Date().getTime() < expiresAt;
  }

  public createAuthHeaderValue(): string {
    return 'Bearer ' + this.accessToken;
  }
}
