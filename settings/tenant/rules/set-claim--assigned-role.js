function SetAssignedRoles(user, context, callback) {
    const namespace = 'https://client.devpie.io/claims/roles';
    const defaultRole = context.idToken[namespace];
    const assignedRoles = (context.authorization || {}).roles;

    let idTokenClaims = context.idToken || {};
    let accessTokenClaims = context.accessToken || {};

    idTokenClaims[namespace] = assignedRoles || defaultRole;
    accessTokenClaims[namespace] = assignedRoles || defaultRole;

    context.idToken = idTokenClaims;
    context.accessToken = accessTokenClaims;

    callback(null, user, context);
}


