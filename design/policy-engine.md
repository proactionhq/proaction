# Policy Engine

This proposal shows how a policy engine could be used to define and execute the Proaction checks, instead of hard coded implementation.

## Goals

- Quicker execution of checks due to less runtime duplicate requests made to the GitHub API
- More testable policy execution because of easy-to-mock datasets
- Ability to quickly create new checks by writing new policy files and not code

## Non Goals

- Changes to any implementation decisions that the policies (in code) today have implemented
- Maximum efficiency and guaranteed no duplication of GitHub requests

## Background

By enabling a pipeline approach that ends with policy execution controlled by a policy engine, it would be easier to create new Proaction checks without a deep understanding of the GitHub API or Proaction codebase.
A policy engine enables a common language to parse inputs, apply rules, and produce outputs.
Using this would allow that a Proaction check could be defined in a single policy file.
This also could eventually be extended to support local (even private) policy files.

To enable this, Proaction would have to support a pipelined execution that is collect -> analyze -> remediate.
The collect step allows the check to request specific information needed as inputs.
Collect then proceeds to collect all of this information from the APIs.
The analyze is the policy execution.
The inputs are provided to each policy and executed.
The remediate step uses the policy output to automatically remediate any possible recommendations.

## High-Level Design

To start, a policy engine needs to be selected.
Two leading choices are OpenPolicyAgent and HashiCorp Sentinel.
A detailed design must be written to allow a check to specify the inputs to pass to the collect phase.
The existing checks would then be written in the policy language.

## Detailed Design

In progress

## Alternatives Considered

In progress

## Security Considerations

In progress
