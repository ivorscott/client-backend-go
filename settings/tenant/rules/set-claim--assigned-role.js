function SetAssignedRoles(user, context, callback) {
    const namespace = 'https://client.devpie.io/claims/roles';
    const defaultRoles = context.idToken[namespace];
    const assignedRoles = (context.authorization || {}).roles.push(defaultRoles)

    let idTokenClaims = context.idToken || {};
    let accessTokenClaims = context.accessToken || {};

    idTokenClaims[namespace] = assignedRoles;
    accessTokenClaims[namespace] = assignedRoles;

    context.idToken = idTokenClaims;
    context.accessToken = accessTokenClaims;

    callback(null, user, context);
}


