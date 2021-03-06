# Copyright 2020 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM alpine:latest as builder
RUN addgroup -S apigee && adduser -S apigee -G apigee -u 1001
RUN apk --update add ca-certificates

FROM scratch
COPY --from=builder /etc/passwd /etc/group /etc/shadow /etc/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY apigeecli / 
USER apigee
ENTRYPOINT ["/apigeecli"]