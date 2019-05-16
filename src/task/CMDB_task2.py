#!/usr/bin/env python
# -*- coding: UTF-8 -*-
import redis

rc = redis.StrictRedis(host='192.168.102.173', port='6379', db=0, password='xZS,JrgQls')
host_info = "{\"beat\":{\"address\":[\"192.168.200.10\",\"192.168.122.1\",\"fe80::fa98:efff:fec6:efd4\"],\"hostname\":\"rbtnode1\",\"name\":\"rbtnode1\",\"version\":\"10.1.7\"},\"bizid\":0,\"cloudid\":0,\"data\":{\"city\":\"\",\"country\":\"\",\"cpu\":{\"cpuinfo\":[{\"cacheSize\":25600,\"coreId\":\"0\",\"cores\":1,\"cpu\":0,\"family\":\"6\",\"flags\":[\"fpu\",\"vme\",\"de\",\"pse\",\"tsc\",\"msr\",\"pae\",\"mce\",\"cx8\",\"apic\",\"sep\",\"mtrr\",\"pge\",\"mca\",\"cmov\",\"pat\",\"pse36\",\"clflush\",\"mmx\",\"fxsr\",\"sse\",\"sse2\",\"ss\",\"syscall\",\"nx\",\"rdtscp\",\"lm\",\"constant_tsc\",\"arch_perfmon\",\"nopl\",\"xtopology\",\"tsc_reliable\",\"nonstop_tsc\",\"pni\",\"pclmulqdq\",\"ssse3\",\"cx16\",\"pcid\",\"sse4_1\",\"sse4_2\",\"x2apic\",\"popcnt\",\"tsc_deadline_timer\",\"aes\",\"xsave\",\"avx\",\"f16c\",\"rdrand\",\"hypervisor\",\"lahf_lm\",\"fsgsbase\",\"tsc_adjust\",\"smep\",\"ibpb\",\"ibrs\",\"stibp\",\"arat\",\"spec_ctrl\",\"intel_stibp\",\"arch_capabilities\"],\"mhz\":2399.999,\"microcode\":\"0x42d\",\"model\":\"62\",\"modelName\":\"Intel(R) Xeon(R) CPU E5-4650 v2 @ 2.40GHz\",\"physicalId\":\"0\",\"stepping\":4,\"vendorId\":\"GenuineIntel\"},{\"cacheSize\":25600,\"coreId\":\"0\",\"cores\":1,\"cpu\":1,\"family\":\"6\",\"flags\":[\"fpu\",\"vme\",\"de\",\"pse\",\"tsc\",\"msr\",\"pae\",\"mce\",\"cx8\",\"apic\",\"sep\",\"mtrr\",\"pge\",\"mca\",\"cmov\",\"pat\",\"pse36\",\"clflush\",\"mmx\",\"fxsr\",\"sse\",\"sse2\",\"ss\",\"syscall\",\"nx\",\"rdtscp\",\"lm\",\"constant_tsc\",\"arch_perfmon\",\"nopl\",\"xtopology\",\"tsc_reliable\",\"nonstop_tsc\",\"pni\",\"pclmulqdq\",\"ssse3\",\"cx16\",\"pcid\",\"sse4_1\",\"sse4_2\",\"x2apic\",\"popcnt\",\"tsc_deadline_timer\",\"aes\",\"xsave\",\"avx\",\"f16c\",\"rdrand\",\"hypervisor\",\"lahf_lm\",\"fsgsbase\",\"tsc_adjust\",\"smep\",\"ibpb\",\"ibrs\",\"stibp\",\"arat\",\"spec_ctrl\",\"intel_stibp\",\"arch_capabilities\"],\"mhz\":2399.999,\"microcode\":\"0x42d\",\"model\":\"62\",\"modelName\":\"Intel(R) Xeon(R) CPU E5-4650 v2 @ 2.40GHz\",\"physicalId\":\"2\",\"stepping\":4,\"vendorId\":\"GenuineIntel\"},{\"cacheSize\":25600,\"coreId\":\"0\",\"cores\":1,\"cpu\":2,\"family\":\"6\",\"flags\":[\"fpu\",\"vme\",\"de\",\"pse\",\"tsc\",\"msr\",\"pae\",\"mce\",\"cx8\",\"apic\",\"sep\",\"mtrr\",\"pge\",\"mca\",\"cmov\",\"pat\",\"pse36\",\"clflush\",\"mmx\",\"fxsr\",\"sse\",\"sse2\",\"ss\",\"syscall\",\"nx\",\"rdtscp\",\"lm\",\"constant_tsc\",\"arch_perfmon\",\"nopl\",\"xtopology\",\"tsc_reliable\",\"nonstop_tsc\",\"pni\",\"pclmulqdq\",\"ssse3\",\"cx16\",\"pcid\",\"sse4_1\",\"sse4_2\",\"x2apic\",\"popcnt\",\"tsc_deadline_timer\",\"aes\",\"xsave\",\"avx\",\"f16c\",\"rdrand\",\"hypervisor\",\"lahf_lm\",\"fsgsbase\",\"tsc_adjust\",\"smep\",\"ibpb\",\"ibrs\",\"stibp\",\"arat\",\"spec_ctrl\",\"intel_stibp\",\"arch_capabilities\"],\"mhz\":2399.999,\"microcode\":\"0x42d\",\"model\":\"62\",\"modelName\":\"Intel(R) Xeon(R) CPU E5-4650 v2 @ 2.40GHz\",\"physicalId\":\"4\",\"stepping\":4,\"vendorId\":\"GenuineIntel\"},{\"cacheSize\":25600,\"coreId\":\"0\",\"cores\":1,\"cpu\":3,\"family\":\"6\",\"flags\":[\"fpu\",\"vme\",\"de\",\"pse\",\"tsc\",\"msr\",\"pae\",\"mce\",\"cx8\",\"apic\",\"sep\",\"mtrr\",\"pge\",\"mca\",\"cmov\",\"pat\",\"pse36\",\"clflush\",\"mmx\",\"fxsr\",\"sse\",\"sse2\",\"ss\",\"syscall\",\"nx\",\"rdtscp\",\"lm\",\"constant_tsc\",\"arch_perfmon\",\"nopl\",\"xtopology\",\"tsc_reliable\",\"nonstop_tsc\",\"pni\",\"pclmulqdq\",\"ssse3\",\"cx16\",\"pcid\",\"sse4_1\",\"sse4_2\",\"x2apic\",\"popcnt\",\"tsc_deadline_timer\",\"aes\",\"xsave\",\"avx\",\"f16c\",\"rdrand\",\"hypervisor\",\"lahf_lm\",\"fsgsbase\",\"tsc_adjust\",\"smep\",\"ibpb\",\"ibrs\",\"stibp\",\"arat\",\"spec_ctrl\",\"intel_stibp\",\"arch_capabilities\"],\"mhz\":2399.999,\"microcode\":\"0x42d\",\"model\":\"62\",\"modelName\":\"Intel(R) Xeon(R) CPU E5-4650 v2 @ 2.40GHz\",\"physicalId\":\"6\",\"stepping\":4,\"vendorId\":\"GenuineIntel\"},{\"cacheSize\":25600,\"coreId\":\"0\",\"cores\":1,\"cpu\":4,\"family\":\"6\",\"flags\":[\"fpu\",\"vme\",\"de\",\"pse\",\"tsc\",\"msr\",\"pae\",\"mce\",\"cx8\",\"apic\",\"sep\",\"mtrr\",\"pge\",\"mca\",\"cmov\",\"pat\",\"pse36\",\"clflush\",\"mmx\",\"fxsr\",\"sse\",\"sse2\",\"ss\",\"syscall\",\"nx\",\"rdtscp\",\"lm\",\"constant_tsc\",\"arch_perfmon\",\"nopl\",\"xtopology\",\"tsc_reliable\",\"nonstop_tsc\",\"pni\",\"pclmulqdq\",\"ssse3\",\"cx16\",\"pcid\",\"sse4_1\",\"sse4_2\",\"x2apic\",\"popcnt\",\"tsc_deadline_timer\",\"aes\",\"xsave\",\"avx\",\"f16c\",\"rdrand\",\"hypervisor\",\"lahf_lm\",\"fsgsbase\",\"tsc_adjust\",\"smep\",\"ibpb\",\"ibrs\",\"stibp\",\"arat\",\"spec_ctrl\",\"intel_stibp\",\"arch_capabilities\"],\"mhz\":2399.999,\"microcode\":\"0x42d\",\"model\":\"62\",\"modelName\":\"Intel(R) Xeon(R) CPU E5-4650 v2 @ 2.40GHz\",\"physicalId\":\"8\",\"stepping\":4,\"vendorId\":\"GenuineIntel\"},{\"cacheSize\":25600,\"coreId\":\"0\",\"cores\":1,\"cpu\":5,\"family\":\"6\",\"flags\":[\"fpu\",\"vme\",\"de\",\"pse\",\"tsc\",\"msr\",\"pae\",\"mce\",\"cx8\",\"apic\",\"sep\",\"mtrr\",\"pge\",\"mca\",\"cmov\",\"pat\",\"pse36\",\"clflush\",\"mmx\",\"fxsr\",\"sse\",\"sse2\",\"ss\",\"syscall\",\"nx\",\"rdtscp\",\"lm\",\"constant_tsc\",\"arch_perfmon\",\"nopl\",\"xtopology\",\"tsc_reliable\",\"nonstop_tsc\",\"pni\",\"pclmulqdq\",\"ssse3\",\"cx16\",\"pcid\",\"sse4_1\",\"sse4_2\",\"x2apic\",\"popcnt\",\"tsc_deadline_timer\",\"aes\",\"xsave\",\"avx\",\"f16c\",\"rdrand\",\"hypervisor\",\"lahf_lm\",\"fsgsbase\",\"tsc_adjust\",\"smep\",\"ibpb\",\"ibrs\",\"stibp\",\"arat\",\"spec_ctrl\",\"intel_stibp\",\"arch_capabilities\"],\"mhz\":2399.999,\"microcode\":\"0x42d\",\"model\":\"62\",\"modelName\":\"Intel(R) Xeon(R) CPU E5-4650 v2 @ 2.40GHz\",\"physicalId\":\"10\",\"stepping\":4,\"vendorId\":\"GenuineIntel\"},{\"cacheSize\":25600,\"coreId\":\"0\",\"cores\":1,\"cpu\":6,\"family\":\"6\",\"flags\":[\"fpu\",\"vme\",\"de\",\"pse\",\"tsc\",\"msr\",\"pae\",\"mce\",\"cx8\",\"apic\",\"sep\",\"mtrr\",\"pge\",\"mca\",\"cmov\",\"pat\",\"pse36\",\"clflush\",\"mmx\",\"fxsr\",\"sse\",\"sse2\",\"ss\",\"syscall\",\"nx\",\"rdtscp\",\"lm\",\"constant_tsc\",\"arch_perfmon\",\"nopl\",\"xtopology\",\"tsc_reliable\",\"nonstop_tsc\",\"pni\",\"pclmulqdq\",\"ssse3\",\"cx16\",\"pcid\",\"sse4_1\",\"sse4_2\",\"x2apic\",\"popcnt\",\"tsc_deadline_timer\",\"aes\",\"xsave\",\"avx\",\"f16c\",\"rdrand\",\"hypervisor\",\"lahf_lm\",\"fsgsbase\",\"tsc_adjust\",\"smep\",\"ibpb\",\"ibrs\",\"stibp\",\"arat\",\"spec_ctrl\",\"intel_stibp\",\"arch_capabilities\"],\"mhz\":2399.999,\"microcode\":\"0x42d\",\"model\":\"62\",\"modelName\":\"Intel(R) Xeon(R) CPU E5-4650 v2 @ 2.40GHz\",\"physicalId\":\"12\",\"stepping\":4,\"vendorId\":\"GenuineIntel\"},{\"cacheSize\":25600,\"coreId\":\"0\",\"cores\":1,\"cpu\":7,\"family\":\"6\",\"flags\":[\"fpu\",\"vme\",\"de\",\"pse\",\"tsc\",\"msr\",\"pae\",\"mce\",\"cx8\",\"apic\",\"sep\",\"mtrr\",\"pge\",\"mca\",\"cmov\",\"pat\",\"pse36\",\"clflush\",\"mmx\",\"fxsr\",\"sse\",\"sse2\",\"ss\",\"syscall\",\"nx\",\"rdtscp\",\"lm\",\"constant_tsc\",\"arch_perfmon\",\"nopl\",\"xtopology\",\"tsc_reliable\",\"nonstop_tsc\",\"pni\",\"pclmulqdq\",\"ssse3\",\"cx16\",\"pcid\",\"sse4_1\",\"sse4_2\",\"x2apic\",\"popcnt\",\"tsc_deadline_timer\",\"aes\",\"xsave\",\"avx\",\"f16c\",\"rdrand\",\"hypervisor\",\"lahf_lm\",\"fsgsbase\",\"tsc_adjust\",\"smep\",\"ibpb\",\"ibrs\",\"stibp\",\"arat\",\"spec_ctrl\",\"intel_stibp\",\"arch_capabilities\"],\"mhz\":2399.999,\"microcode\":\"0x42d\",\"model\":\"62\",\"modelName\":\"Intel(R) Xeon(R) CPU E5-4650 v2 @ 2.40GHz\",\"physicalId\":\"14\",\"stepping\":4,\"vendorId\":\"GenuineIntel\"},{\"cacheSize\":25600,\"coreId\":\"0\",\"cores\":1,\"cpu\":8,\"family\":\"6\",\"flags\":[\"fpu\",\"vme\",\"de\",\"pse\",\"tsc\",\"msr\",\"pae\",\"mce\",\"cx8\",\"apic\",\"sep\",\"mtrr\",\"pge\",\"mca\",\"cmov\",\"pat\",\"pse36\",\"clflush\",\"mmx\",\"fxsr\",\"sse\",\"sse2\",\"ss\",\"syscall\",\"nx\",\"rdtscp\",\"lm\",\"constant_tsc\",\"arch_perfmon\",\"nopl\",\"xtopology\",\"tsc_reliable\",\"nonstop_tsc\",\"pni\",\"pclmulqdq\",\"ssse3\",\"cx16\",\"pcid\",\"sse4_1\",\"sse4_2\",\"x2apic\",\"popcnt\",\"tsc_deadline_timer\",\"aes\",\"xsave\",\"avx\",\"f16c\",\"rdrand\",\"hypervisor\",\"lahf_lm\",\"fsgsbase\",\"tsc_adjust\",\"smep\",\"ibpb\",\"ibrs\",\"stibp\",\"arat\",\"spec_ctrl\",\"intel_stibp\",\"arch_capabilities\"],\"mhz\":2399.999,\"microcode\":\"0x42d\",\"model\":\"62\",\"modelName\":\"Intel(R) Xeon(R) CPU E5-4650 v2 @ 2.40GHz\",\"physicalId\":\"16\",\"stepping\":4,\"vendorId\":\"GenuineIntel\"},{\"cacheSize\":25600,\"coreId\":\"0\",\"cores\":1,\"cpu\":9,\"family\":\"6\",\"flags\":[\"fpu\",\"vme\",\"de\",\"pse\",\"tsc\",\"msr\",\"pae\",\"mce\",\"cx8\",\"apic\",\"sep\",\"mtrr\",\"pge\",\"mca\",\"cmov\",\"pat\",\"pse36\",\"clflush\",\"mmx\",\"fxsr\",\"sse\",\"sse2\",\"ss\",\"syscall\",\"nx\",\"rdtscp\",\"lm\",\"constant_tsc\",\"arch_perfmon\",\"nopl\",\"xtopology\",\"tsc_reliable\",\"nonstop_tsc\",\"pni\",\"pclmulqdq\",\"ssse3\",\"cx16\",\"pcid\",\"sse4_1\",\"sse4_2\",\"x2apic\",\"popcnt\",\"tsc_deadline_timer\",\"aes\",\"xsave\",\"avx\",\"f16c\",\"rdrand\",\"hypervisor\",\"lahf_lm\",\"fsgsbase\",\"tsc_adjust\",\"smep\",\"ibpb\",\"ibrs\",\"stibp\",\"arat\",\"spec_ctrl\",\"intel_stibp\",\"arch_capabilities\"],\"mhz\":2399.999,\"microcode\":\"0x42d\",\"model\":\"62\",\"modelName\":\"Intel(R) Xeon(R) CPU E5-4650 v2 @ 2.40GHz\",\"physicalId\":\"18\",\"stepping\":4,\"vendorId\":\"GenuineIntel\"},{\"cacheSize\":25600,\"coreId\":\"0\",\"cores\":1,\"cpu\":10,\"family\":\"6\",\"flags\":[\"fpu\",\"vme\",\"de\",\"pse\",\"tsc\",\"msr\",\"pae\",\"mce\",\"cx8\",\"apic\",\"sep\",\"mtrr\",\"pge\",\"mca\",\"cmov\",\"pat\",\"pse36\",\"clflush\",\"mmx\",\"fxsr\",\"sse\",\"sse2\",\"ss\",\"syscall\",\"nx\",\"rdtscp\",\"lm\",\"constant_tsc\",\"arch_perfmon\",\"nopl\",\"xtopology\",\"tsc_reliable\",\"nonstop_tsc\",\"pni\",\"pclmulqdq\",\"ssse3\",\"cx16\",\"pcid\",\"sse4_1\",\"sse4_2\",\"x2apic\",\"popcnt\",\"tsc_deadline_timer\",\"aes\",\"xsave\",\"avx\",\"f16c\",\"rdrand\",\"hypervisor\",\"lahf_lm\",\"fsgsbase\",\"tsc_adjust\",\"smep\",\"ibpb\",\"ibrs\",\"stibp\",\"arat\",\"spec_ctrl\",\"intel_stibp\",\"arch_capabilities\"],\"mhz\":2399.999,\"microcode\":\"0x42d\",\"model\":\"62\",\"modelName\":\"Intel(R) Xeon(R) CPU E5-4650 v2 @ 2.40GHz\",\"physicalId\":\"20\",\"stepping\":4,\"vendorId\":\"GenuineIntel\"},{\"cacheSize\":25600,\"coreId\":\"0\",\"cores\":1,\"cpu\":11,\"family\":\"6\",\"flags\":[\"fpu\",\"vme\",\"de\",\"pse\",\"tsc\",\"msr\",\"pae\",\"mce\",\"cx8\",\"apic\",\"sep\",\"mtrr\",\"pge\",\"mca\",\"cmov\",\"pat\",\"pse36\",\"clflush\",\"mmx\",\"fxsr\",\"sse\",\"sse2\",\"ss\",\"syscall\",\"nx\",\"rdtscp\",\"lm\",\"constant_tsc\",\"arch_perfmon\",\"nopl\",\"xtopology\",\"tsc_reliable\",\"nonstop_tsc\",\"pni\",\"pclmulqdq\",\"ssse3\",\"cx16\",\"pcid\",\"sse4_1\",\"sse4_2\",\"x2apic\",\"popcnt\",\"tsc_deadline_timer\",\"aes\",\"xsave\",\"avx\",\"f16c\",\"rdrand\",\"hypervisor\",\"lahf_lm\",\"fsgsbase\",\"tsc_adjust\",\"smep\",\"ibpb\",\"ibrs\",\"stibp\",\"arat\",\"spec_ctrl\",\"intel_stibp\",\"arch_capabilities\"],\"mhz\":2399.999,\"microcode\":\"0x42d\",\"model\":\"62\",\"modelName\":\"Intel(R) Xeon(R) CPU E5-4650 v2 @ 2.40GHz\",\"physicalId\":\"22\",\"stepping\":4,\"vendorId\":\"GenuineIntel\"},{\"cacheSize\":25600,\"coreId\":\"0\",\"cores\":1,\"cpu\":12,\"family\":\"6\",\"flags\":[\"fpu\",\"vme\",\"de\",\"pse\",\"tsc\",\"msr\",\"pae\",\"mce\",\"cx8\",\"apic\",\"sep\",\"mtrr\",\"pge\",\"mca\",\"cmov\",\"pat\",\"pse36\",\"clflush\",\"mmx\",\"fxsr\",\"sse\",\"sse2\",\"ss\",\"syscall\",\"nx\",\"rdtscp\",\"lm\",\"constant_tsc\",\"arch_perfmon\",\"nopl\",\"xtopology\",\"tsc_reliable\",\"nonstop_tsc\",\"pni\",\"pclmulqdq\",\"ssse3\",\"cx16\",\"pcid\",\"sse4_1\",\"sse4_2\",\"x2apic\",\"popcnt\",\"tsc_deadline_timer\",\"aes\",\"xsave\",\"avx\",\"f16c\",\"rdrand\",\"hypervisor\",\"lahf_lm\",\"fsgsbase\",\"tsc_adjust\",\"smep\",\"ibpb\",\"ibrs\",\"stibp\",\"arat\",\"spec_ctrl\",\"intel_stibp\",\"arch_capabilities\"],\"mhz\":2399.999,\"microcode\":\"0x42d\",\"model\":\"62\",\"modelName\":\"Intel(R) Xeon(R) CPU E5-4650 v2 @ 2.40GHz\",\"physicalId\":\"24\",\"stepping\":4,\"vendorId\":\"GenuineIntel\"},{\"cacheSize\":25600,\"coreId\":\"0\",\"cores\":1,\"cpu\":13,\"family\":\"6\",\"flags\":[\"fpu\",\"vme\",\"de\",\"pse\",\"tsc\",\"msr\",\"pae\",\"mce\",\"cx8\",\"apic\",\"sep\",\"mtrr\",\"pge\",\"mca\",\"cmov\",\"pat\",\"pse36\",\"clflush\",\"mmx\",\"fxsr\",\"sse\",\"sse2\",\"ss\",\"syscall\",\"nx\",\"rdtscp\",\"lm\",\"constant_tsc\",\"arch_perfmon\",\"nopl\",\"xtopology\",\"tsc_reliable\",\"nonstop_tsc\",\"pni\",\"pclmulqdq\",\"ssse3\",\"cx16\",\"pcid\",\"sse4_1\",\"sse4_2\",\"x2apic\",\"popcnt\",\"tsc_deadline_timer\",\"aes\",\"xsave\",\"avx\",\"f16c\",\"rdrand\",\"hypervisor\",\"lahf_lm\",\"fsgsbase\",\"tsc_adjust\",\"smep\",\"ibpb\",\"ibrs\",\"stibp\",\"arat\",\"spec_ctrl\",\"intel_stibp\",\"arch_capabilities\"],\"mhz\":2399.999,\"microcode\":\"0x42d\",\"model\":\"62\",\"modelName\":\"Intel(R) Xeon(R) CPU E5-4650 v2 @ 2.40GHz\",\"physicalId\":\"26\",\"stepping\":4,\"vendorId\":\"GenuineIntel\"},{\"cacheSize\":25600,\"coreId\":\"0\",\"cores\":1,\"cpu\":14,\"family\":\"6\",\"flags\":[\"fpu\",\"vme\",\"de\",\"pse\",\"tsc\",\"msr\",\"pae\",\"mce\",\"cx8\",\"apic\",\"sep\",\"mtrr\",\"pge\",\"mca\",\"cmov\",\"pat\",\"pse36\",\"clflush\",\"mmx\",\"fxsr\",\"sse\",\"sse2\",\"ss\",\"syscall\",\"nx\",\"rdtscp\",\"lm\",\"constant_tsc\",\"arch_perfmon\",\"nopl\",\"xtopology\",\"tsc_reliable\",\"nonstop_tsc\",\"pni\",\"pclmulqdq\",\"ssse3\",\"cx16\",\"pcid\",\"sse4_1\",\"sse4_2\",\"x2apic\",\"popcnt\",\"tsc_deadline_timer\",\"aes\",\"xsave\",\"avx\",\"f16c\",\"rdrand\",\"hypervisor\",\"lahf_lm\",\"fsgsbase\",\"tsc_adjust\",\"smep\",\"ibpb\",\"ibrs\",\"stibp\",\"arat\",\"spec_ctrl\",\"intel_stibp\",\"arch_capabilities\"],\"mhz\":2399.999,\"microcode\":\"0x42d\",\"model\":\"62\",\"modelName\":\"Intel(R) Xeon(R) CPU E5-4650 v2 @ 2.40GHz\",\"physicalId\":\"28\",\"stepping\":4,\"vendorId\":\"GenuineIntel\"},{\"cacheSize\":25600,\"coreId\":\"0\",\"cores\":1,\"cpu\":15,\"family\":\"6\",\"flags\":[\"fpu\",\"vme\",\"de\",\"pse\",\"tsc\",\"msr\",\"pae\",\"mce\",\"cx8\",\"apic\",\"sep\",\"mtrr\",\"pge\",\"mca\",\"cmov\",\"pat\",\"pse36\",\"clflush\",\"mmx\",\"fxsr\",\"sse\",\"sse2\",\"ss\",\"syscall\",\"nx\",\"rdtscp\",\"lm\",\"constant_tsc\",\"arch_perfmon\",\"nopl\",\"xtopology\",\"tsc_reliable\",\"nonstop_tsc\",\"pni\",\"pclmulqdq\",\"ssse3\",\"cx16\",\"pcid\",\"sse4_1\",\"sse4_2\",\"x2apic\",\"popcnt\",\"tsc_deadline_timer\",\"aes\",\"xsave\",\"avx\",\"f16c\",\"rdrand\",\"hypervisor\",\"lahf_lm\",\"fsgsbase\",\"tsc_adjust\",\"smep\",\"ibpb\",\"ibrs\",\"stibp\",\"arat\",\"spec_ctrl\",\"intel_stibp\",\"arch_capabilities\"],\"mhz\":2399.999,\"microcode\":\"0x42d\",\"model\":\"62\",\"modelName\":\"Intel(R) Xeon(R) CPU E5-4650 v2 @ 2.40GHz\",\"physicalId\":\"30\",\"stepping\":4,\"vendorId\":\"GenuineIntel\"}],\"per_stat\":[{\"cpu\":\"cpu0\",\"guest\":0,\"guestNice\":0,\"idle\":821397.4,\"iowait\":4901.79,\"irq\":0,\"nice\":2.78,\"softirq\":252197.95,\"steal\":0,\"stolen\":0,\"system\":673479.88,\"user\":1917770.54},{\"cpu\":\"cpu1\",\"guest\":0,\"guestNice\":0,\"idle\":697986.86,\"iowait\":4749.44,\"irq\":0,\"nice\":3.28,\"softirq\":18998.01,\"steal\":0,\"stolen\":0,\"system\":761140.53,\"user\":2523477.79},{\"cpu\":\"cpu2\",\"guest\":0,\"guestNice\":0,\"idle\":687345.02,\"iowait\":4719.98,\"irq\":0,\"nice\":2.81,\"softirq\":17248.4,\"steal\":0,\"stolen\":0,\"system\":763103.52,\"user\":2527868.45},{\"cpu\":\"cpu3\",\"guest\":0,\"guestNice\":0,\"idle\":681394.63,\"iowait\":4714.69,\"irq\":0,\"nice\":3.14,\"softirq\":16551.7,\"steal\":0,\"stolen\":0,\"system\":763196.41,\"user\":2527900.27},{\"cpu\":\"cpu4\",\"guest\":0,\"guestNice\":0,\"idle\":678982.08,\"iowait\":3608.11,\"irq\":0,\"nice\":4.03,\"softirq\":16602.6,\"steal\":0,\"stolen\":0,\"system\":759815.03,\"user\":2513493.58},{\"cpu\":\"cpu5\",\"guest\":0,\"guestNice\":0,\"idle\":773230.19,\"iowait\":6238.73,\"irq\":0,\"nice\":1.75,\"softirq\":284321.87,\"steal\":0,\"stolen\":0,\"system\":672545.81,\"user\":1890652.05},{\"cpu\":\"cpu6\",\"guest\":0,\"guestNice\":0,\"idle\":774492.2,\"iowait\":4890.33,\"irq\":0,\"nice\":2.39,\"softirq\":228138.73,\"steal\":0,\"stolen\":0,\"system\":686907.01,\"user\":1973997.25},{\"cpu\":\"cpu7\",\"guest\":0,\"guestNice\":0,\"idle\":785432.74,\"iowait\":5528.03,\"irq\":0,\"nice\":2.63,\"softirq\":251248.64,\"steal\":0,\"stolen\":0,\"system\":684491.13,\"user\":1917290.19},{\"cpu\":\"cpu8\",\"guest\":0,\"guestNice\":0,\"idle\":786020.76,\"iowait\":5847.98,\"irq\":0,\"nice\":3,\"softirq\":244467.74,\"steal\":0,\"stolen\":0,\"system\":668063.25,\"user\":1960364.79},{\"cpu\":\"cpu9\",\"guest\":0,\"guestNice\":0,\"idle\":658255.84,\"iowait\":4857.57,\"irq\":0,\"nice\":5.49,\"softirq\":18305.49,\"steal\":0,\"stolen\":0,\"system\":760895.23,\"user\":2549160.52},{\"cpu\":\"cpu10\",\"guest\":0,\"guestNice\":0,\"idle\":650304.05,\"iowait\":4885.38,\"irq\":0,\"nice\":3.71,\"softirq\":14973.53,\"steal\":0,\"stolen\":0,\"system\":761515.55,\"user\":2554153.5},{\"cpu\":\"cpu11\",\"guest\":0,\"guestNice\":0,\"idle\":638376.81,\"iowait\":4846.26,\"irq\":0,\"nice\":3.89,\"softirq\":14731.71,\"steal\":0,\"stolen\":0,\"system\":759911.03,\"user\":2560156.54},{\"cpu\":\"cpu12\",\"guest\":0,\"guestNice\":0,\"idle\":631136.23,\"iowait\":4790.81,\"irq\":0,\"nice\":3.12,\"softirq\":14796.73,\"steal\":0,\"stolen\":0,\"system\":756769.11,\"user\":2563534.49},{\"cpu\":\"cpu13\",\"guest\":0,\"guestNice\":0,\"idle\":731817.87,\"iowait\":5380.38,\"irq\":0,\"nice\":3.19,\"softirq\":225374.25,\"steal\":0,\"stolen\":0,\"system\":688496.08,\"user\":2006078.22},{\"cpu\":\"cpu14\",\"guest\":0,\"guestNice\":0,\"idle\":726855.44,\"iowait\":4985.89,\"irq\":0,\"nice\":2.89,\"softirq\":292715.56,\"steal\":0,\"stolen\":0,\"system\":661904.14,\"user\":1920076.19},{\"cpu\":\"cpu15\",\"guest\":0,\"guestNice\":0,\"idle\":705139.8,\"iowait\":5741.67,\"irq\":0,\"nice\":2.54,\"softirq\":235392.05,\"steal\":0,\"stolen\":0,\"system\":675273.09,\"user\":2011095.85}],\"per_usage\":[83.48946136209368,82.11306765651624,73.1751824821705,95.36082474523751,82.06364512794882,77.68691588746917,89.21800947743094,75.00000000134119,81.17386489746205,96.14624506246678,86.93957114672669,90.69306930799887,92.92035398481313,89.61625281710278,86.66666666523936,93.1764705818804],\"total_stat\":{\"cpu\":\"cpu-total\",\"guest\":0,\"guestNice\":0,\"idle\":11428167.98,\"iowait\":80687.11,\"irq\":0,\"nice\":50.71,\"softirq\":2146065.03,\"steal\":0,\"stolen\":0,\"system\":11497506.87,\"user\":35917070.28},\"total_usage\":85.9505766060639},\"datetime\":\"2019-05-16 14:59:17\",\"disk\":{\"diskstat\":{\"dm-0\":{\"avgqu_sz\":0.009054216458531095,\"avgrq_sz\":17.393700787401574,\"await\":4.283464566929134,\"ioTime\":19630381,\"iopsInProgress\":0,\"major\":253,\"mergedReadCount\":0,\"mergedWriteCount\":0,\"minor\":0,\"name\":\"dm-0\",\"readBytes\":3080555520,\"readCount\":86642,\"readSectors\":6016710,\"readTime\":4383680,\"serialNumber\":\"\",\"speedByteRead\":0,\"speedByteWrite\":18824.248618254296,\"speedIORead\":0,\"speedIOWrite\":2.113760092340899,\"svctm\":2.0078740157480315,\"util\":0.0042441639649364515,\"weightedIoTime\":59574393,\"writeBytes\":113230268928,\"writeCount\":12823421,\"writeSectors\":221152869,\"writeTime\":55117434},\"dm-1\":{\"avgqu_sz\":0,\"avgrq_sz\":0,\"await\":0,\"ioTime\":181389,\"iopsInProgress\":0,\"major\":253,\"mergedReadCount\":0,\"mergedWriteCount\":0,\"minor\":1,\"name\":\"dm-1\",\"readBytes\":13373440,\"readCount\":1664,\"readSectors\":26120,\"readTime\":16249,\"serialNumber\":\"\",\"speedByteRead\":0,\"speedByteWrite\":0,\"speedIORead\":0,\"speedIOWrite\":0,\"svctm\":0,\"util\":0,\"weightedIoTime\":460807,\"writeBytes\":60653568,\"writeCount\":14808,\"writeSectors\":118464,\"writeTime\":444554},\"dm-2\":{\"avgqu_sz\":0,\"avgrq_sz\":0,\"await\":0,\"ioTime\":164328,\"iopsInProgress\":0,\"major\":253,\"mergedReadCount\":0,\"mergedWriteCount\":0,\"minor\":2,\"name\":\"dm-2\",\"readBytes\":8611328,\"readCount\":271,\"readSectors\":16819,\"readTime\":1567,\"serialNumber\":\"\",\"speedByteRead\":0,\"speedByteWrite\":0,\"speedIORead\":0,\"speedIOWrite\":0,\"svctm\":0,\"util\":0,\"weightedIoTime\":273125,\"writeBytes\":357848576,\"writeCount\":77375,\"writeSectors\":698923,\"writeTime\":271517},\"dm-3\":{\"avgqu_sz\":0.20385302055898688,\"avgrq_sz\":21.890319258496394,\"await\":2.352986611740474,\"ioTime\":454636865,\"iopsInProgress\":0,\"major\":253,\"mergedReadCount\":0,\"mergedWriteCount\":0,\"minor\":3,\"name\":\"dm-3\",\"readBytes\":181045675008,\"readCount\":1073002,\"readSectors\":353604834,\"readTime\":79851073,\"serialNumber\":\"\",\"speedByteRead\":0,\"speedByteWrite\":724524.7922232761,\"speedIORead\":0,\"speedIOWrite\":64.64444250907128,\"svctm\":1.392378990731205,\"util\":0.09000956361716207,\"weightedIoTime\":1075943893,\"writeBytes\":2545550489600,\"writeCount\":270475933,\"writeSectors\":4971778300,\"writeTime\":993883398},\"sda\":{\"avgqu_sz\":0.1521075077472715,\"avgrq_sz\":21.042790456431536,\"await\":2.37551867219917,\"ioTime\":458938909,\"iopsInProgress\":0,\"major\":8,\"mergedReadCount\":2898,\"mergedWriteCount\":2534004,\"minor\":0,\"name\":\"sda\",\"readBytes\":104152297472,\"readCount\":949350,\"readSectors\":203422456,\"readTime\":50927748,\"serialNumber\":\"\",\"speedByteRead\":0,\"speedByteWrite\":691452.4025051027,\"speedIORead\":0,\"speedIOWrite\":64.178416661941,\"svctm\":1.4299792531120332,\"util\":0.09177380432415526,\"weightedIoTime\":881226849,\"writeBytes\":2339086463488,\"writeCount\":271112374,\"writeSectors\":4568528249,\"writeTime\":833447270},\"sda1\":{\"avgqu_sz\":0,\"avgrq_sz\":0,\"await\":0,\"ioTime\":157601,\"iopsInProgress\":0,\"major\":8,\"mergedReadCount\":0,\"mergedWriteCount\":0,\"minor\":1,\"name\":\"sda1\",\"readBytes\":26782720,\"readCount\":532,\"readSectors\":52310,\"readTime\":10773,\"serialNumber\":\"\",\"speedByteRead\":0,\"speedByteWrite\":0,\"speedIORead\":0,\"speedIOWrite\":0,\"svctm\":0,\"util\":0,\"weightedIoTime\":272197,\"writeBytes\":254970368,\"writeCount\":75334,\"writeSectors\":497989,\"writeTime\":261700},\"sda2\":{\"avgqu_sz\":0.008804559754711304,\"avgrq_sz\":19.04310344827586,\"await\":4.568965517241379,\"ioTime\":19854909,\"iopsInProgress\":0,\"major\":8,\"mergedReadCount\":1317,\"mergedWriteCount\":768325,\"minor\":2,\"name\":\"sda2\",\"readBytes\":3105096192,\"readCount\":87133,\"readSectors\":6064641,\"readTime\":4340982,\"serialNumber\":\"\",\"speedByteRead\":0,\"speedByteWrite\":18824.248618254296,\"speedIORead\":0,\"speedIOWrite\":1.930678509539719,\"svctm\":2.189655172413793,\"util\":0.004227520184681798,\"weightedIoTime\":55349139,\"writeBytes\":113648820224,\"writeCount\":12147291,\"writeSectors\":221970352,\"writeTime\":51324877},\"sda3\":{\"avgqu_sz\":0.1433029479925602,\"avgrq_sz\":21.1048128342246,\"await\":2.307486631016043,\"ioTime\":443729044,\"iopsInProgress\":0,\"major\":8,\"mergedReadCount\":1581,\"mergedWriteCount\":1765679,\"minor\":3,\"name\":\"sda3\",\"readBytes\":101015020032,\"readCount\":861503,\"readSectors\":197294961,\"readTime\":46574988,\"serialNumber\":\"\",\"speedByteRead\":0,\"speedByteWrite\":672628.1538868485,\"speedIORead\":0,\"speedIOWrite\":62.24773815240128,\"svctm\":1.4318181818181819,\"util\":0.08912744326366547,\"weightedIoTime\":825599515,\"writeBytes\":2225182672896,\"writeCount\":258889749,\"writeSectors\":4346059908,\"writeTime\":781860693},\"sdb\":{\"avgqu_sz\":0.004393957987228326,\"avgrq_sz\":20.435294117647057,\"await\":3.1058823529411765,\"ioTime\":19317193,\"iopsInProgress\":0,\"major\":8,\"mergedReadCount\":5,\"mergedWriteCount\":218127,\"minor\":16,\"name\":\"sdb\",\"readBytes\":79977774592,\"readCount\":208640,\"readSectors\":156206591,\"readTime\":33175230,\"serialNumber\":\"\",\"speedByteRead\":0,\"speedByteWrite\":14802.046106793896,\"speedIORead\":0,\"speedIOWrite\":1.4147213216454837,\"svctm\":2.5647058823529414,\"util\":0.0036283440955142995,\"weightedIoTime\":101752296,\"writeBytes\":192841715200,\"writeCount\":9215274,\"writeSectors\":376643975,\"writeTime\":68673307},\"sdc\":{\"avgqu_sz\":0.0018308158280118025,\"avgrq_sz\":136.03125,\"await\":3.4375,\"ioTime\":1552998,\"iopsInProgress\":0,\"major\":8,\"mergedReadCount\":0,\"mergedWriteCount\":4027,\"minor\":32,\"name\":\"sdc\",\"readBytes\":70652928,\"readCount\":1620,\"readSectors\":137994,\"readTime\":52278,\"serialNumber\":\"\",\"speedByteRead\":0,\"speedByteWrite\":37094.59222963375,\"speedIORead\":0,\"speedIOWrite\":0.532600968148888,\"svctm\":2.8125,\"util\":0.0014979402229187476,\"weightedIoTime\":129427941,\"writeBytes\":127526277632,\"writeCount\":383120,\"writeSectors\":249074761,\"writeTime\":129376557},\"sdc1\":{\"avgqu_sz\":0.0018308158280118025,\"avgrq_sz\":136.03125,\"await\":3.4375,\"ioTime\":1552218,\"iopsInProgress\":0,\"major\":8,\"mergedReadCount\":0,\"mergedWriteCount\":4027,\"minor\":33,\"name\":\"sdc1\",\"readBytes\":65545216,\"readCount\":1419,\"readSectors\":128018,\"readTime\":51976,\"serialNumber\":\"\",\"speedByteRead\":0,\"speedByteWrite\":37094.59222963375,\"speedIORead\":0,\"speedIOWrite\":0.532600968148888,\"svctm\":2.8125,\"util\":0.0014979402229187476,\"weightedIoTime\":129427113,\"writeBytes\":127526273536,\"writeCount\":383119,\"writeSectors\":249074753,\"writeTime\":129376031},\"sr0\":{\"avgqu_sz\":0,\"avgrq_sz\":0,\"await\":0,\"ioTime\":19,\"iopsInProgress\":0,\"major\":11,\"mergedReadCount\":0,\"mergedWriteCount\":0,\"minor\":0,\"name\":\"sr0\",\"readBytes\":1282048,\"readCount\":24,\"readSectors\":2504,\"readTime\":22,\"serialNumber\":\"VMware_Virtual_SATA_CDRW_Drive_00000000000000000001\",\"speedByteRead\":0,\"speedByteWrite\":0,\"speedIORead\":0,\"speedIOWrite\":0,\"svctm\":0,\"util\":0,\"weightedIoTime\":22,\"writeBytes\":0,\"writeCount\":0,\"writeSectors\":0,\"writeTime\":0}},\"partition\":[{\"device\":\"/dev/mapper/centos-root\",\"fstype\":\"xfs\",\"mountpoint\":\"/\",\"opts\":\"rw,relatime,attr2,inode64,noquota\"},{\"device\":\"/dev/sda1\",\"fstype\":\"xfs\",\"mountpoint\":\"/boot\",\"opts\":\"rw,relatime,attr2,inode64,noquota\"},{\"device\":\"/dev/mapper/centos-data\",\"fstype\":\"xfs\",\"mountpoint\":\"/data\",\"opts\":\"rw,relatime,attr2,inode64,noquota\"},{\"device\":\"/dev/mapper/centos-home\",\"fstype\":\"xfs\",\"mountpoint\":\"/home\",\"opts\":\"rw,relatime,attr2,inode64,noquota\"},{\"device\":\"/dev/mapper/centos-data\",\"fstype\":\"xfs\",\"mountpoint\":\"/data/bkee/public/paas_agent/docker/overlay\",\"opts\":\"rw,relatime,attr2,inode64,noquota\"}],\"usage\":[{\"free\":13603770368,\"fstype\":\"xfs\",\"inodesFree\":10337764,\"inodesTotal\":10485760,\"inodesUsed\":147996,\"inodesUsedPercent\":1.4113998413085938,\"path\":\"/\",\"total\":21464350720,\"used\":7860580352,\"usedPercent\":36.62156127870054},{\"free\":886112256,\"fstype\":\"xfs\",\"inodesFree\":523961,\"inodesTotal\":524288,\"inodesUsed\":327,\"inodesUsedPercent\":0.06237030029296875,\"path\":\"/boot\",\"total\":1063256064,\"used\":177143808,\"usedPercent\":16.660502958579883},{\"free\":100340154368,\"fstype\":\"xfs\",\"inodesFree\":195979987,\"inodesTotal\":196367328,\"inodesUsed\":387341,\"inodesUsedPercent\":0.19725328237903203,\"path\":\"/data\",\"total\":582876676096,\"used\":482536521728,\"usedPercent\":82.78535434972973},{\"free\":21426479104,\"fstype\":\"xfs\",\"inodesFree\":10485596,\"inodesTotal\":10485760,\"inodesUsed\":164,\"inodesUsedPercent\":0.0015640258789062498,\"path\":\"/home\",\"total\":21464350720,\"used\":37871616,\"usedPercent\":0.17643960674157305},{\"free\":100340154368,\"fstype\":\"xfs\",\"inodesFree\":195979987,\"inodesTotal\":196367328,\"inodesUsed\":387341,\"inodesUsedPercent\":0.19725328237903203,\"path\":\"/data/bkee/public/paas_agent/docker/overlay\",\"total\":582876676096,\"used\":482536521728,\"usedPercent\":82.78535434972973}]},\"env\":{\"crontab\":[{\"content\":\"\\n\\n\\n\\n\\n\\n\\n\\n*/10 * * * * /data/bkee/bkdata/dataapi/bin/update_cc_cache.sh\\n\\n# \xe5\xae\x9a\xe6\x97\xb6\xe5\x88\xa0\xe9\x99\xa4ES\xe4\xb8\xad\xe7\x9a\x84\xe8\xbf\x87\xe6\x9c\x9f\xe6\x95\xb0\xe6\x8d\xae\\n0 */1 * * * /data/bkee/bkdata/dataapi/datatools/scripts/es_delete_expired_index.sh\\n\\n# \xe5\xae\x9a\xe6\x97\xb6\xe6\x8f\x90\xe5\x89\x8d\xe5\x88\x9b\xe5\xbb\xbaES\xe4\xb8\xad\xe7\x9a\x84\xe7\xb4\xa2\xe5\xbc\x95\\n*/10 * * * * /data/bkee/bkdata/dataapi/datatools/scripts/es_create_mapping_ahead.sh\\n# \xe5\xae\x9a\xe6\x97\xb6\xe5\x88\xa0\xe9\x99\xa4Crate\xe4\xb8\xad\xe7\x9a\x84\xe8\xbf\x87\xe6\x9c\x9f\xe6\x95\xb0\xe6\x8d\xae\\n0 */1 * * * /data/bkee/bkdata/dataapi/datatools/scripts/crate_tool.sh\\n* * * * * export INSTALL_PATH=/data/bkee; /data/bkee/bin/process_watch rabbitmq \\u003e/dev/null 2\\u003e\\u00261\\n* * * * * export INSTALL_PATH=/data/bkee; /data/bkee/bin/process_watch beanstalk \\u003e/dev/null 2\\u003e\\u00261\\n* * * * * export INSTALL_PATH=/data/bkee; /data/bkee/bin/process_watch etcd \\u003e/dev/null 2\\u003e\\u00261\\n* * * * * export INSTALL_PATH=/data/bkee; /data/bkee/bin/process_watch license \\u003e/dev/null 2\\u003e\\u00261\\n* * * * * export INSTALL_PATH=/data/bkee; /data/bkee/bin/process_watch redis \\u003e/dev/null 2\\u003e\\u00261\\n* * * * * /usr/local/gse/agent/bin/gsectl watch\\n* * * * * export INSTALL_PATH=/data/bkee; /data/bkee/bin/process_watch paas_agent \\u003e/dev/null 2\\u003e\\u00261\\n* * * * * export INSTALL_PATH=/data/bkee; /data/bkee/bin/process_watch saas-o \\u003e/dev/null 2\\u003e\\u00261\\n* * * * * export INSTALL_PATH=/data/bkee; /data/bkee/bin/process_watch consul \\u003e/dev/null 2\\u003e\\u00261\\n* * * * * export INSTALL_PATH=/data/bkee; /data/bkee/bin/process_watch bkdata \\u003e/dev/null 2\\u003e\\u00261\\n* * * * * export INSTALL_PATH=/data/bkee; /data/bkee/bin/process_watch mysql \\u003e/dev/null 2\\u003e\\u00261\\n\",\"user\":\"root\"}],\"host\":\"127.0.0.1   localhost localhost.localdomain localhost4 localhost4.localdomain4\\n::1         localhost localhost.localdomain localhost6 localhost6.localdomain6\\n192.168.102.173 paas.bk.com\\n192.168.102.173 cmdb.bk.com\\n192.168.102.173 job.bk.com\\n192.168.102.174 paas.bk.com\\n192.168.102.174 cmdb.bk.com\\n192.168.102.174 job.bk.com\\n192.168.200.10   gse-1\\n192.168.102.175   rbtnode2\\n\\n192.168.163.155 gitlab.canway.net\\n192.168.200.10   rbtnode1\\n192.168.102.173 paas.sgis.gov.cn\\n192.168.102.173 cmdb.sgis.gov.cn\\n192.168.102.173 job.sgis.gov.cn\\n192.168.102.174 paas.sgis.gov.cn\\n192.168.102.174 cmdb.sgis.gov.cn\\n192.168.102.174 job.sgis.gov.cn\\n\",\"route\":\"Kernel IP routing table\\nDestination     Gateway         Genmask         Flags Metric Ref    Use Iface\\n0.0.0.0         192.168.102.1   0.0.0.0         UG    0      0        0 ens192\\n169.254.0.0     0.0.0.0         255.255.0.0     U     1002   0        0 ens192\\n192.168.102.0   0.0.0.0         255.255.255.0   U     0      0        0 ens192\\n192.168.122.0   0.0.0.0         255.255.255.0   U     0      0        0 virbr0\\n\"},\"load\":{\"load_avg\":{\"load1\":44.36,\"load15\":49.22,\"load5\":47.91}},\"mem\":{\"meminfo\":{\"active\":28414009344,\"available\":10939748352,\"buffers\":0,\"cached\":12132409344,\"dirty\":1359872,\"free\":610230272,\"inactive\":10699468800,\"total\":42020630528,\"used\":31080882176,\"usedPercent\":73.96576820828422,\"wired\":0,\"writeback\":0,\"writebacktmp\":0},\"vmstat\":{\"free\":17148407808,\"sin\":5664768,\"sout\":60653568,\"total\":17175670784,\"used\":27262976,\"usedPercent\":0.15873019658362822}},\"net\":{\"dev\":[{\"bytesRecv\":4277678393666,\"bytesSent\":3280563353355,\"dropin\":37073,\"dropout\":0,\"errin\":0,\"errout\":0,\"fifoin\":0,\"fifoout\":0,\"name\":\"ens192\",\"packetsRecv\":13141674951,\"packetsSent\":9206822395,\"speedPacketsRecv\":2321,\"speedPacketsSent\":1504,\"speedRecv\":844543,\"speedSent\":411800},{\"bytesRecv\":554124901500,\"bytesSent\":554124901500,\"dropin\":0,\"dropout\":0,\"errin\":0,\"errout\":0,\"fifoin\":0,\"fifoout\":0,\"name\":\"lo\",\"packetsRecv\":2245387493,\"packetsSent\":2245387493,\"speedPacketsRecv\":469,\"speedPacketsSent\":469,\"speedRecv\":131851,\"speedSent\":131851},{\"bytesRecv\":0,\"bytesSent\":0,\"dropin\":0,\"dropout\":0,\"errin\":0,\"errout\":0,\"fifoin\":0,\"fifoout\":0,\"name\":\"virbr0-nic\",\"packetsRecv\":0,\"packetsSent\":0,\"speedPacketsRecv\":0,\"speedPacketsSent\":0,\"speedRecv\":0,\"speedSent\":0},{\"bytesRecv\":0,\"bytesSent\":0,\"dropin\":0,\"dropout\":0,\"errin\":0,\"errout\":0,\"fifoin\":0,\"fifoout\":0,\"name\":\"virbr0\",\"packetsRecv\":0,\"packetsSent\":0,\"speedPacketsRecv\":0,\"speedPacketsSent\":0,\"speedRecv\":0,\"speedSent\":0}],\"interface\":[{\"addrs\":[{\"addr\":\"127.0.0.1/8\"},{\"addr\":\"::1/128\"}],\"flags\":[\"up\",\"loopback\"],\"hardwareaddr\":\"\",\"mtu\":65536,\"name\":\"lo\"},{\"addrs\":[{\"addr\":\"192.168.200.10/24\"},{\"addr\":\"fe80::fa98:efff:fec6:efd4/64\"}],\"flags\":[\"up\",\"broadcast\",\"multicast\"],\"hardwareaddr\":\"f8:98:ef:c6:ef:d4\",\"mtu\":1500,\"name\":\"ens192\"},{\"addrs\":[{\"addr\":\"192.168.122.1/24\"}],\"flags\":[\"up\",\"broadcast\",\"multicast\"],\"hardwareaddr\":\"52:54:00:34:f1:df\",\"mtu\":1500,\"name\":\"virbr0\"},{\"addrs\":[],\"flags\":[\"broadcast\",\"multicast\"],\"hardwareaddr\":\"52:54:00:34:f1:df\",\"mtu\":1500,\"name\":\"virbr0-nic\"}],\"netstat\":{\"close\":0,\"closeWait\":53,\"closing\":0,\"established\":2427,\"finWait1\":0,\"finWait2\":0,\"lastAck\":0,\"listen\":57,\"synRecv\":2,\"syncSent\":0,\"timeWait\":711},\"protocolstat\":{\"udp\":{\"inCsumErrors\":0,\"inDatagrams\":1898,\"inErrors\":0,\"noPorts\":0,\"outDatagrams\":1888,\"rcvbufErrors\":0,\"sndbufErrors\":0}}},\"system\":{\"docker\":{\"Client\":{\"ApiVersion\":\"1.24\",\"GoVersion\":\"go1.6.4\",\"Version\":\"1.12.4\"},\"Server\":{\"ApiVersion\":\"1.24\",\"GoVersion\":\"go1.6.4\",\"Version\":\"1.12.4\"}},\"info\":{\"bootTime\":1552639310,\"hostid\":\"B9342F42-85A8-E67A-B85B-2B68B679BFF9\",\"hostname\":\"rbtnode1\",\"kernelVersion\":\"3.10.0-862.el7.x86_64\",\"os\":\"linux\",\"platform\":\"centos\",\"platformFamily\":\"rhel\",\"platformVersion\":\"7.5.1804\",\"procs\":994,\"systemtype\":\"64-bit\",\"uptime\":5350602,\"virtualizationRole\":\"\",\"virtualizationSystem\":\"\"}},\"timezone\":8,\"utctime\":\"2019-05-16 06:59:17\"},\"dataid\":1001,\"gseindex\":83509,\"ip\":\"192.168.200.10\",\"type\":\"basereport\"}"

rc.publish("snapshot2", str(host_info))
