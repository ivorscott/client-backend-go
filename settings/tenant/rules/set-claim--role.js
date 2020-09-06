function SetRole(user, context, callback) {
    const namespace = 'https://client.devpie.io/claims/';
    if (user.email.indexOf('@devpie.io') != -1) {
        context.idToken[namespace + 'roles'] = 'admin';
    } else {
        context.idToken[namespace + 'roles'] = 'user';
    }
    callback(null, user, context);
}