
# S2PD: A naive reverse proxy for steps.dodgson.org

S2PD is a HTTP reverse proxy specifically implemented for handling
pages on dodgson.org, which are statically generated and
hosted at Amazon S3 backend.

Why do I need this? I have some old links to redirect, and I don't want to learn
DSL for proper reverse proxy servers. Also, I was looking for a justification to
write some server in Go language and this perfectly fits to the desire.
