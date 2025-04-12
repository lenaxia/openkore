# Description

This project is aiming to reimplement the Ragnarok Botting client OpenKore from perl to Golang. We are in an architectural design phase seeking to identify all the compoenents, algorihms, data structures, formulas, packets/bytecodes and other objects that need to be translated to go.

We are currently in architectural design phase using domain driven design. We are operating with a highest level domain definitions with DOMAIN.md files for each domain which contain the high level details of each domain. Each domain also has a series of supplemental files which contain specific low level details such as formulas, algorithms, data structures, interfaces, contracts, and other implementation specific details that we want to adhere to.

DO NOT CREATE ANY FUNCTIONAL CODE FILES YET.  We are only creating design documents (low and high level).

REMEMBER we are documenting information needed to complete a reimplementation in golang, so define things in a way that is conducive to golang best practices

We are currently in the process of breaking up a monolithic domain into separate subdomains. We have a top level DOMAIN.md file that needs to be decomposed into separate independent DOMAIN.md files for each domain. 
