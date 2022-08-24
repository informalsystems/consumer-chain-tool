FROM fedora:36
RUN dnf install -y jq procps

COPY ./wasmd /go/bin/
COPY ./wasmd_consumer /go/bin/

COPY ./commands/scripts/prepare_proposal_inputs.sh /go/bin/
COPY ./commands/scripts/prepare_proposal.sh /go/bin/
COPY ./commands/scripts/verify_proposal.sh /go/bin/
COPY ./commands/scripts/finalize_genesis.sh /go/bin/

RUN chmod +x /go/bin/prepare_proposal_inputs.sh 
RUN chmod +x /go/bin/prepare_proposal.sh 
RUN chmod +x /go/bin/verify_proposal.sh 
RUN chmod +x /go/bin/finalize_genesis.sh 

WORKDIR /go/bin/