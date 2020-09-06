function SetDefaultRoles(user, context, callback) {
    const namespace = 'https://client.devpie.io/claims/roles';

    if (user.email.indexOf('@devpie.io') !== -1) {
        context.idToken[namespace] = 'employee';
    } else {
        context.idToken[namespace] = 'user';
    }

    callback(null, user, context);
}