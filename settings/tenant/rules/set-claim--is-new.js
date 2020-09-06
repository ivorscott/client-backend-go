function IsNew (user, context, callback) {
    const namespace = 'devpie-client';
    context.idToken[namespace + 'is_new'] = context.stats.loginsCount === 1;
    callback(null, user, context);
}