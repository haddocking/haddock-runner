package input

// This is a very ugly hack to allow for "repeated" fields in the YAML file.
type ModuleParams struct {
	Order []string `yaml:"order"`
	// Analysis
	Alascan           map[string]any `yaml:"alascan"`
	Alascan_0         map[string]any `yaml:"alascan.0"`
	Alascan_1         map[string]any `yaml:"alascan.1"`
	Alascan_2         map[string]any `yaml:"alascan.2"`
	Alascan_3         map[string]any `yaml:"alascan.3"`
	Alascan_4         map[string]any `yaml:"alascan.4"`
	Alascan_5         map[string]any `yaml:"alascan.5"`
	Alascan_6         map[string]any `yaml:"alascan.6"`
	Alascan_7         map[string]any `yaml:"alascan.7"`
	Alascan_8         map[string]any `yaml:"alascan.8"`
	Alascan_9         map[string]any `yaml:"alascan.9"`
	Alascan_10        map[string]any `yaml:"alascan.10"`
	Alascan_11        map[string]any `yaml:"alascan.11"`
	Alascan_12        map[string]any `yaml:"alascan.12"`
	Alascan_13        map[string]any `yaml:"alascan.13"`
	Alascan_14        map[string]any `yaml:"alascan.14"`
	Alascan_15        map[string]any `yaml:"alascan.15"`
	Alascan_16        map[string]any `yaml:"alascan.16"`
	Alascan_17        map[string]any `yaml:"alascan.17"`
	Alascan_18        map[string]any `yaml:"alascan.18"`
	Alascan_19        map[string]any `yaml:"alascan.19"`
	Alascan_20        map[string]any `yaml:"alascan.20"`
	Alascan_21        map[string]any `yaml:"alascan.21"`
	Alascan_22        map[string]any `yaml:"alascan.22"`
	Alascan_23        map[string]any `yaml:"alascan.23"`
	Alascan_24        map[string]any `yaml:"alascan.24"`
	Alascan_25        map[string]any `yaml:"alascan.25"`
	Caprieval         map[string]any `yaml:"caprieval"`
	Caprieval_0       map[string]any `yaml:"caprieval.0"`
	Caprieval_1       map[string]any `yaml:"caprieval.1"`
	Caprieval_2       map[string]any `yaml:"caprieval.2"`
	Caprieval_3       map[string]any `yaml:"caprieval.3"`
	Caprieval_4       map[string]any `yaml:"caprieval.4"`
	Caprieval_5       map[string]any `yaml:"caprieval.5"`
	Caprieval_6       map[string]any `yaml:"caprieval.6"`
	Caprieval_7       map[string]any `yaml:"caprieval.7"`
	Caprieval_8       map[string]any `yaml:"caprieval.8"`
	Caprieval_9       map[string]any `yaml:"caprieval.9"`
	Caprieval_10      map[string]any `yaml:"caprieval.10"`
	Caprieval_11      map[string]any `yaml:"caprieval.11"`
	Caprieval_12      map[string]any `yaml:"caprieval.12"`
	Caprieval_13      map[string]any `yaml:"caprieval.13"`
	Caprieval_14      map[string]any `yaml:"caprieval.14"`
	Caprieval_15      map[string]any `yaml:"caprieval.15"`
	Caprieval_16      map[string]any `yaml:"caprieval.16"`
	Caprieval_17      map[string]any `yaml:"caprieval.17"`
	Caprieval_18      map[string]any `yaml:"caprieval.18"`
	Caprieval_19      map[string]any `yaml:"caprieval.19"`
	Caprieval_20      map[string]any `yaml:"caprieval.20"`
	Caprieval_21      map[string]any `yaml:"caprieval.21"`
	Caprieval_22      map[string]any `yaml:"caprieval.22"`
	Caprieval_23      map[string]any `yaml:"caprieval.23"`
	Caprieval_24      map[string]any `yaml:"caprieval.24"`
	Caprieval_25      map[string]any `yaml:"caprieval.25"`
	Clustfcc          map[string]any `yaml:"clustfcc"`
	Clustfcc_0        map[string]any `yaml:"clustfcc.0"`
	Clustfcc_1        map[string]any `yaml:"clustfcc.1"`
	Clustfcc_2        map[string]any `yaml:"clustfcc.2"`
	Clustfcc_3        map[string]any `yaml:"clustfcc.3"`
	Clustfcc_4        map[string]any `yaml:"clustfcc.4"`
	Clustfcc_5        map[string]any `yaml:"clustfcc.5"`
	Clustfcc_6        map[string]any `yaml:"clustfcc.6"`
	Clustfcc_7        map[string]any `yaml:"clustfcc.7"`
	Clustfcc_8        map[string]any `yaml:"clustfcc.8"`
	Clustfcc_9        map[string]any `yaml:"clustfcc.9"`
	Clustfcc_10       map[string]any `yaml:"clustfcc.10"`
	Clustfcc_11       map[string]any `yaml:"clustfcc.11"`
	Clustfcc_12       map[string]any `yaml:"clustfcc.12"`
	Clustfcc_13       map[string]any `yaml:"clustfcc.13"`
	Clustfcc_14       map[string]any `yaml:"clustfcc.14"`
	Clustfcc_15       map[string]any `yaml:"clustfcc.15"`
	Clustfcc_16       map[string]any `yaml:"clustfcc.16"`
	Clustfcc_17       map[string]any `yaml:"clustfcc.17"`
	Clustfcc_18       map[string]any `yaml:"clustfcc.18"`
	Clustfcc_19       map[string]any `yaml:"clustfcc.19"`
	Clustfcc_20       map[string]any `yaml:"clustfcc.20"`
	Clustfcc_21       map[string]any `yaml:"clustfcc.21"`
	Clustfcc_22       map[string]any `yaml:"clustfcc.22"`
	Clustfcc_23       map[string]any `yaml:"clustfcc.23"`
	Clustfcc_24       map[string]any `yaml:"clustfcc.24"`
	Clustfcc_25       map[string]any `yaml:"clustfcc.25"`
	Clustrmsd         map[string]any `yaml:"clustrmsd"`
	Clustrmsd_0       map[string]any `yaml:"clustrmsd.0"`
	Clustrmsd_1       map[string]any `yaml:"clustrmsd.1"`
	Clustrmsd_2       map[string]any `yaml:"clustrmsd.2"`
	Clustrmsd_3       map[string]any `yaml:"clustrmsd.3"`
	Clustrmsd_4       map[string]any `yaml:"clustrmsd.4"`
	Clustrmsd_5       map[string]any `yaml:"clustrmsd.5"`
	Clustrmsd_6       map[string]any `yaml:"clustrmsd.6"`
	Clustrmsd_7       map[string]any `yaml:"clustrmsd.7"`
	Clustrmsd_8       map[string]any `yaml:"clustrmsd.8"`
	Clustrmsd_9       map[string]any `yaml:"clustrmsd.9"`
	Clustrmsd_10      map[string]any `yaml:"clustrmsd.10"`
	Clustrmsd_11      map[string]any `yaml:"clustrmsd.11"`
	Clustrmsd_12      map[string]any `yaml:"clustrmsd.12"`
	Clustrmsd_13      map[string]any `yaml:"clustrmsd.13"`
	Clustrmsd_14      map[string]any `yaml:"clustrmsd.14"`
	Clustrmsd_15      map[string]any `yaml:"clustrmsd.15"`
	Clustrmsd_16      map[string]any `yaml:"clustrmsd.16"`
	Clustrmsd_17      map[string]any `yaml:"clustrmsd.17"`
	Clustrmsd_18      map[string]any `yaml:"clustrmsd.18"`
	Clustrmsd_19      map[string]any `yaml:"clustrmsd.19"`
	Clustrmsd_20      map[string]any `yaml:"clustrmsd.20"`
	Clustrmsd_21      map[string]any `yaml:"clustrmsd.21"`
	Clustrmsd_22      map[string]any `yaml:"clustrmsd.22"`
	Clustrmsd_23      map[string]any `yaml:"clustrmsd.23"`
	Clustrmsd_24      map[string]any `yaml:"clustrmsd.24"`
	Clustrmsd_25      map[string]any `yaml:"clustrmsd.25"`
	Contactmap        map[string]any `yaml:"contactmap"`
	Contactmap_0      map[string]any `yaml:"contactmap.0"`
	Contactmap_1      map[string]any `yaml:"contactmap.1"`
	Contactmap_2      map[string]any `yaml:"contactmap.2"`
	Contactmap_3      map[string]any `yaml:"contactmap.3"`
	Contactmap_4      map[string]any `yaml:"contactmap.4"`
	Contactmap_5      map[string]any `yaml:"contactmap.5"`
	Contactmap_6      map[string]any `yaml:"contactmap.6"`
	Contactmap_7      map[string]any `yaml:"contactmap.7"`
	Contactmap_8      map[string]any `yaml:"contactmap.8"`
	Contactmap_9      map[string]any `yaml:"contactmap.9"`
	Contactmap_10     map[string]any `yaml:"contactmap.10"`
	Contactmap_11     map[string]any `yaml:"contactmap.11"`
	Contactmap_12     map[string]any `yaml:"contactmap.12"`
	Contactmap_13     map[string]any `yaml:"contactmap.13"`
	Contactmap_14     map[string]any `yaml:"contactmap.14"`
	Contactmap_15     map[string]any `yaml:"contactmap.15"`
	Contactmap_16     map[string]any `yaml:"contactmap.16"`
	Contactmap_17     map[string]any `yaml:"contactmap.17"`
	Contactmap_18     map[string]any `yaml:"contactmap.18"`
	Contactmap_19     map[string]any `yaml:"contactmap.19"`
	Contactmap_20     map[string]any `yaml:"contactmap.20"`
	Contactmap_21     map[string]any `yaml:"contactmap.21"`
	Contactmap_22     map[string]any `yaml:"contactmap.22"`
	Contactmap_23     map[string]any `yaml:"contactmap.23"`
	Contactmap_24     map[string]any `yaml:"contactmap.24"`
	Contactmap_25     map[string]any `yaml:"contactmap.25"`
	Ilrmsdmatrix      map[string]any `yaml:"ilrmsdmatrix"`
	Ilrmsdmatrix_0    map[string]any `yaml:"ilrmsdmatrix.0"`
	Ilrmsdmatrix_1    map[string]any `yaml:"ilrmsdmatrix.1"`
	Ilrmsdmatrix_2    map[string]any `yaml:"ilrmsdmatrix.2"`
	Ilrmsdmatrix_3    map[string]any `yaml:"ilrmsdmatrix.3"`
	Ilrmsdmatrix_4    map[string]any `yaml:"ilrmsdmatrix.4"`
	Ilrmsdmatrix_5    map[string]any `yaml:"ilrmsdmatrix.5"`
	Ilrmsdmatrix_6    map[string]any `yaml:"ilrmsdmatrix.6"`
	Ilrmsdmatrix_7    map[string]any `yaml:"ilrmsdmatrix.7"`
	Ilrmsdmatrix_8    map[string]any `yaml:"ilrmsdmatrix.8"`
	Ilrmsdmatrix_9    map[string]any `yaml:"ilrmsdmatrix.9"`
	Ilrmsdmatrix_10   map[string]any `yaml:"ilrmsdmatrix.10"`
	Ilrmsdmatrix_11   map[string]any `yaml:"ilrmsdmatrix.11"`
	Ilrmsdmatrix_12   map[string]any `yaml:"ilrmsdmatrix.12"`
	Ilrmsdmatrix_13   map[string]any `yaml:"ilrmsdmatrix.13"`
	Ilrmsdmatrix_14   map[string]any `yaml:"ilrmsdmatrix.14"`
	Ilrmsdmatrix_15   map[string]any `yaml:"ilrmsdmatrix.15"`
	Ilrmsdmatrix_16   map[string]any `yaml:"ilrmsdmatrix.16"`
	Ilrmsdmatrix_17   map[string]any `yaml:"ilrmsdmatrix.17"`
	Ilrmsdmatrix_18   map[string]any `yaml:"ilrmsdmatrix.18"`
	Ilrmsdmatrix_19   map[string]any `yaml:"ilrmsdmatrix.19"`
	Ilrmsdmatrix_20   map[string]any `yaml:"ilrmsdmatrix.20"`
	Ilrmsdmatrix_21   map[string]any `yaml:"ilrmsdmatrix.21"`
	Ilrmsdmatrix_22   map[string]any `yaml:"ilrmsdmatrix.22"`
	Ilrmsdmatrix_23   map[string]any `yaml:"ilrmsdmatrix.23"`
	Ilrmsdmatrix_24   map[string]any `yaml:"ilrmsdmatrix.24"`
	Ilrmsdmatrix_25   map[string]any `yaml:"ilrmsdmatrix.25"`
	Filter            map[string]any `yaml:"filter"`
	Filter_0          map[string]any `yaml:"filter.0"`
	Filter_1          map[string]any `yaml:"filter.1"`
	Filter_2          map[string]any `yaml:"filter.2"`
	Filter_3          map[string]any `yaml:"filter.3"`
	Filter_4          map[string]any `yaml:"filter.4"`
	Filter_5          map[string]any `yaml:"filter.5"`
	Filter_6          map[string]any `yaml:"filter.6"`
	Filter_7          map[string]any `yaml:"filter.7"`
	Filter_8          map[string]any `yaml:"filter.8"`
	Filter_9          map[string]any `yaml:"filter.9"`
	Filter_10         map[string]any `yaml:"filter.10"`
	Filter_11         map[string]any `yaml:"filter.11"`
	Filter_12         map[string]any `yaml:"filter.12"`
	Filter_13         map[string]any `yaml:"filter.13"`
	Filter_14         map[string]any `yaml:"filter.14"`
	Filter_15         map[string]any `yaml:"filter.15"`
	Filter_16         map[string]any `yaml:"filter.16"`
	Filter_17         map[string]any `yaml:"filter.17"`
	Filter_18         map[string]any `yaml:"filter.18"`
	Filter_19         map[string]any `yaml:"filter.19"`
	Filter_20         map[string]any `yaml:"filter.20"`
	Filter_21         map[string]any `yaml:"filter.21"`
	Filter_22         map[string]any `yaml:"filter.22"`
	Filter_23         map[string]any `yaml:"filter.23"`
	Filter_24         map[string]any `yaml:"filter.24"`
	Filter_25         map[string]any `yaml:"filter.25"`
	Rmsdmatrix        map[string]any `yaml:"rmsdmatrix"`
	Rmsdmatrix_0      map[string]any `yaml:"rmsdmatrix.0"`
	Rmsdmatrix_1      map[string]any `yaml:"rmsdmatrix.1"`
	Rmsdmatrix_2      map[string]any `yaml:"rmsdmatrix.2"`
	Rmsdmatrix_3      map[string]any `yaml:"rmsdmatrix.3"`
	Rmsdmatrix_4      map[string]any `yaml:"rmsdmatrix.4"`
	Rmsdmatrix_5      map[string]any `yaml:"rmsdmatrix.5"`
	Rmsdmatrix_6      map[string]any `yaml:"rmsdmatrix.6"`
	Rmsdmatrix_7      map[string]any `yaml:"rmsdmatrix.7"`
	Rmsdmatrix_8      map[string]any `yaml:"rmsdmatrix.8"`
	Rmsdmatrix_9      map[string]any `yaml:"rmsdmatrix.9"`
	Rmsdmatrix_10     map[string]any `yaml:"rmsdmatrix.10"`
	Rmsdmatrix_11     map[string]any `yaml:"rmsdmatrix.11"`
	Rmsdmatrix_12     map[string]any `yaml:"rmsdmatrix.12"`
	Rmsdmatrix_13     map[string]any `yaml:"rmsdmatrix.13"`
	Rmsdmatrix_14     map[string]any `yaml:"rmsdmatrix.14"`
	Rmsdmatrix_15     map[string]any `yaml:"rmsdmatrix.15"`
	Rmsdmatrix_16     map[string]any `yaml:"rmsdmatrix.16"`
	Rmsdmatrix_17     map[string]any `yaml:"rmsdmatrix.17"`
	Rmsdmatrix_18     map[string]any `yaml:"rmsdmatrix.18"`
	Rmsdmatrix_19     map[string]any `yaml:"rmsdmatrix.19"`
	Rmsdmatrix_20     map[string]any `yaml:"rmsdmatrix.20"`
	Rmsdmatrix_21     map[string]any `yaml:"rmsdmatrix.21"`
	Rmsdmatrix_22     map[string]any `yaml:"rmsdmatrix.22"`
	Rmsdmatrix_23     map[string]any `yaml:"rmsdmatrix.23"`
	Rmsdmatrix_24     map[string]any `yaml:"rmsdmatrix.24"`
	Rmsdmatrix_25     map[string]any `yaml:"rmsdmatrix.25"`
	Seletop           map[string]any `yaml:"seletop"`
	Seletop_0         map[string]any `yaml:"seletop.0"`
	Seletop_1         map[string]any `yaml:"seletop.1"`
	Seletop_2         map[string]any `yaml:"seletop.2"`
	Seletop_3         map[string]any `yaml:"seletop.3"`
	Seletop_4         map[string]any `yaml:"seletop.4"`
	Seletop_5         map[string]any `yaml:"seletop.5"`
	Seletop_6         map[string]any `yaml:"seletop.6"`
	Seletop_7         map[string]any `yaml:"seletop.7"`
	Seletop_8         map[string]any `yaml:"seletop.8"`
	Seletop_9         map[string]any `yaml:"seletop.9"`
	Seletop_10        map[string]any `yaml:"seletop.10"`
	Seletop_11        map[string]any `yaml:"seletop.11"`
	Seletop_12        map[string]any `yaml:"seletop.12"`
	Seletop_13        map[string]any `yaml:"seletop.13"`
	Seletop_14        map[string]any `yaml:"seletop.14"`
	Seletop_15        map[string]any `yaml:"seletop.15"`
	Seletop_16        map[string]any `yaml:"seletop.16"`
	Seletop_17        map[string]any `yaml:"seletop.17"`
	Seletop_18        map[string]any `yaml:"seletop.18"`
	Seletop_19        map[string]any `yaml:"seletop.19"`
	Seletop_20        map[string]any `yaml:"seletop.20"`
	Seletop_21        map[string]any `yaml:"seletop.21"`
	Seletop_22        map[string]any `yaml:"seletop.22"`
	Seletop_23        map[string]any `yaml:"seletop.23"`
	Seletop_24        map[string]any `yaml:"seletop.24"`
	Seletop_25        map[string]any `yaml:"seletop.25"`
	Seletopclusts     map[string]any `yaml:"seletopclusts"`
	Seletopclusts_0   map[string]any `yaml:"seletopclusts.0"`
	Seletopclusts_1   map[string]any `yaml:"seletopclusts.1"`
	Seletopclusts_2   map[string]any `yaml:"seletopclusts.2"`
	Seletopclusts_3   map[string]any `yaml:"seletopclusts.3"`
	Seletopclusts_4   map[string]any `yaml:"seletopclusts.4"`
	Seletopclusts_5   map[string]any `yaml:"seletopclusts.5"`
	Seletopclusts_6   map[string]any `yaml:"seletopclusts.6"`
	Seletopclusts_7   map[string]any `yaml:"seletopclusts.7"`
	Seletopclusts_8   map[string]any `yaml:"seletopclusts.8"`
	Seletopclusts_9   map[string]any `yaml:"seletopclusts.9"`
	Seletopclusts_10  map[string]any `yaml:"seletopclusts.10"`
	Seletopclusts_11  map[string]any `yaml:"seletopclusts.11"`
	Seletopclusts_12  map[string]any `yaml:"seletopclusts.12"`
	Seletopclusts_13  map[string]any `yaml:"seletopclusts.13"`
	Seletopclusts_14  map[string]any `yaml:"seletopclusts.14"`
	Seletopclusts_15  map[string]any `yaml:"seletopclusts.15"`
	Seletopclusts_16  map[string]any `yaml:"seletopclusts.16"`
	Seletopclusts_17  map[string]any `yaml:"seletopclusts.17"`
	Seletopclusts_18  map[string]any `yaml:"seletopclusts.18"`
	Seletopclusts_19  map[string]any `yaml:"seletopclusts.19"`
	Seletopclusts_20  map[string]any `yaml:"seletopclusts.20"`
	Seletopclusts_21  map[string]any `yaml:"seletopclusts.21"`
	Seletopclusts_22  map[string]any `yaml:"seletopclusts.22"`
	Seletopclusts_23  map[string]any `yaml:"seletopclusts.23"`
	Seletopclusts_24  map[string]any `yaml:"seletopclusts.24"`
	Seletopclusts_25  map[string]any `yaml:"seletopclusts.25"`
	Cgtoaa            map[string]any `yaml:"cgtoaa"`
	Cgtoaa_0          map[string]any `yaml:"cgtoaa.0"`
	Cgtoaa_1          map[string]any `yaml:"cgtoaa.1"`
	Cgtoaa_2          map[string]any `yaml:"cgtoaa.2"`
	Cgtoaa_3          map[string]any `yaml:"cgtoaa.3"`
	Cgtoaa_4          map[string]any `yaml:"cgtoaa.4"`
	Cgtoaa_5          map[string]any `yaml:"cgtoaa.5"`
	Cgtoaa_6          map[string]any `yaml:"cgtoaa.6"`
	Cgtoaa_7          map[string]any `yaml:"cgtoaa.7"`
	Cgtoaa_8          map[string]any `yaml:"cgtoaa.8"`
	Cgtoaa_9          map[string]any `yaml:"cgtoaa.9"`
	Cgtoaa_10         map[string]any `yaml:"cgtoaa.10"`
	Cgtoaa_11         map[string]any `yaml:"cgtoaa.11"`
	Cgtoaa_12         map[string]any `yaml:"cgtoaa.12"`
	Cgtoaa_13         map[string]any `yaml:"cgtoaa.13"`
	Cgtoaa_14         map[string]any `yaml:"cgtoaa.14"`
	Cgtoaa_15         map[string]any `yaml:"cgtoaa.15"`
	Cgtoaa_16         map[string]any `yaml:"cgtoaa.16"`
	Cgtoaa_17         map[string]any `yaml:"cgtoaa.17"`
	Cgtoaa_18         map[string]any `yaml:"cgtoaa.18"`
	Cgtoaa_19         map[string]any `yaml:"cgtoaa.19"`
	Cgtoaa_20         map[string]any `yaml:"cgtoaa.20"`
	Cgtoaa_21         map[string]any `yaml:"cgtoaa.21"`
	Cgtoaa_22         map[string]any `yaml:"cgtoaa.22"`
	Cgtoaa_23         map[string]any `yaml:"cgtoaa.23"`
	Cgtoaa_24         map[string]any `yaml:"cgtoaa.24"`
	Cgtoaa_25         map[string]any `yaml:"cgtoaa.25"`
	Emref             map[string]any `yaml:"emref"`
	Emref_0           map[string]any `yaml:"emref.0"`
	Emref_1           map[string]any `yaml:"emref.1"`
	Emref_2           map[string]any `yaml:"emref.2"`
	Emref_3           map[string]any `yaml:"emref.3"`
	Emref_4           map[string]any `yaml:"emref.4"`
	Emref_5           map[string]any `yaml:"emref.5"`
	Emref_6           map[string]any `yaml:"emref.6"`
	Emref_7           map[string]any `yaml:"emref.7"`
	Emref_8           map[string]any `yaml:"emref.8"`
	Emref_9           map[string]any `yaml:"emref.9"`
	Emref_10          map[string]any `yaml:"emref.10"`
	Emref_11          map[string]any `yaml:"emref.11"`
	Emref_12          map[string]any `yaml:"emref.12"`
	Emref_13          map[string]any `yaml:"emref.13"`
	Emref_14          map[string]any `yaml:"emref.14"`
	Emref_15          map[string]any `yaml:"emref.15"`
	Emref_16          map[string]any `yaml:"emref.16"`
	Emref_17          map[string]any `yaml:"emref.17"`
	Emref_18          map[string]any `yaml:"emref.18"`
	Emref_19          map[string]any `yaml:"emref.19"`
	Emref_20          map[string]any `yaml:"emref.20"`
	Emref_21          map[string]any `yaml:"emref.21"`
	Emref_22          map[string]any `yaml:"emref.22"`
	Emref_23          map[string]any `yaml:"emref.23"`
	Emref_24          map[string]any `yaml:"emref.24"`
	Emref_25          map[string]any `yaml:"emref.25"`
	Flexref           map[string]any `yaml:"flexref"`
	Flexref_0         map[string]any `yaml:"flexref.0"`
	Flexref_1         map[string]any `yaml:"flexref.1"`
	Flexref_2         map[string]any `yaml:"flexref.2"`
	Flexref_3         map[string]any `yaml:"flexref.3"`
	Flexref_4         map[string]any `yaml:"flexref.4"`
	Flexref_5         map[string]any `yaml:"flexref.5"`
	Flexref_6         map[string]any `yaml:"flexref.6"`
	Flexref_7         map[string]any `yaml:"flexref.7"`
	Flexref_8         map[string]any `yaml:"flexref.8"`
	Flexref_9         map[string]any `yaml:"flexref.9"`
	Flexref_10        map[string]any `yaml:"flexref.10"`
	Flexref_11        map[string]any `yaml:"flexref.11"`
	Flexref_12        map[string]any `yaml:"flexref.12"`
	Flexref_13        map[string]any `yaml:"flexref.13"`
	Flexref_14        map[string]any `yaml:"flexref.14"`
	Flexref_15        map[string]any `yaml:"flexref.15"`
	Flexref_16        map[string]any `yaml:"flexref.16"`
	Flexref_17        map[string]any `yaml:"flexref.17"`
	Flexref_18        map[string]any `yaml:"flexref.18"`
	Flexref_19        map[string]any `yaml:"flexref.19"`
	Flexref_20        map[string]any `yaml:"flexref.20"`
	Flexref_21        map[string]any `yaml:"flexref.21"`
	Flexref_22        map[string]any `yaml:"flexref.22"`
	Flexref_23        map[string]any `yaml:"flexref.23"`
	Flexref_24        map[string]any `yaml:"flexref.24"`
	Flexref_25        map[string]any `yaml:"flexref.25"`
	Mdref             map[string]any `yaml:"mdref"`
	Mdref_0           map[string]any `yaml:"mdref.0"`
	Mdref_1           map[string]any `yaml:"mdref.1"`
	Mdref_2           map[string]any `yaml:"mdref.2"`
	Mdref_3           map[string]any `yaml:"mdref.3"`
	Mdref_4           map[string]any `yaml:"mdref.4"`
	Mdref_5           map[string]any `yaml:"mdref.5"`
	Mdref_6           map[string]any `yaml:"mdref.6"`
	Mdref_7           map[string]any `yaml:"mdref.7"`
	Mdref_8           map[string]any `yaml:"mdref.8"`
	Mdref_9           map[string]any `yaml:"mdref.9"`
	Mdref_10          map[string]any `yaml:"mdref.10"`
	Mdref_11          map[string]any `yaml:"mdref.11"`
	Mdref_12          map[string]any `yaml:"mdref.12"`
	Mdref_13          map[string]any `yaml:"mdref.13"`
	Mdref_14          map[string]any `yaml:"mdref.14"`
	Mdref_15          map[string]any `yaml:"mdref.15"`
	Mdref_16          map[string]any `yaml:"mdref.16"`
	Mdref_17          map[string]any `yaml:"mdref.17"`
	Mdref_18          map[string]any `yaml:"mdref.18"`
	Mdref_19          map[string]any `yaml:"mdref.19"`
	Mdref_20          map[string]any `yaml:"mdref.20"`
	Mdref_21          map[string]any `yaml:"mdref.21"`
	Mdref_22          map[string]any `yaml:"mdref.22"`
	Mdref_23          map[string]any `yaml:"mdref.23"`
	Mdref_24          map[string]any `yaml:"mdref.24"`
	Mdref_25          map[string]any `yaml:"mdref.25"`
	Openmm            map[string]any `yaml:"openmm"`
	Openmm_0          map[string]any `yaml:"openmm.0"`
	Openmm_1          map[string]any `yaml:"openmm.1"`
	Openmm_2          map[string]any `yaml:"openmm.2"`
	Openmm_3          map[string]any `yaml:"openmm.3"`
	Openmm_4          map[string]any `yaml:"openmm.4"`
	Openmm_5          map[string]any `yaml:"openmm.5"`
	Openmm_6          map[string]any `yaml:"openmm.6"`
	Openmm_7          map[string]any `yaml:"openmm.7"`
	Openmm_8          map[string]any `yaml:"openmm.8"`
	Openmm_9          map[string]any `yaml:"openmm.9"`
	Openmm_10         map[string]any `yaml:"openmm.10"`
	Openmm_11         map[string]any `yaml:"openmm.11"`
	Openmm_12         map[string]any `yaml:"openmm.12"`
	Openmm_13         map[string]any `yaml:"openmm.13"`
	Openmm_14         map[string]any `yaml:"openmm.14"`
	Openmm_15         map[string]any `yaml:"openmm.15"`
	Openmm_16         map[string]any `yaml:"openmm.16"`
	Openmm_17         map[string]any `yaml:"openmm.17"`
	Openmm_18         map[string]any `yaml:"openmm.18"`
	Openmm_19         map[string]any `yaml:"openmm.19"`
	Openmm_20         map[string]any `yaml:"openmm.20"`
	Openmm_21         map[string]any `yaml:"openmm.21"`
	Openmm_22         map[string]any `yaml:"openmm.22"`
	Openmm_23         map[string]any `yaml:"openmm.23"`
	Openmm_24         map[string]any `yaml:"openmm.24"`
	Openmm_25         map[string]any `yaml:"openmm.25"`
	Gdock             map[string]any `yaml:"gdock"`
	Gdock_0           map[string]any `yaml:"gdock.0"`
	Gdock_1           map[string]any `yaml:"gdock.1"`
	Gdock_2           map[string]any `yaml:"gdock.2"`
	Gdock_3           map[string]any `yaml:"gdock.3"`
	Gdock_4           map[string]any `yaml:"gdock.4"`
	Gdock_5           map[string]any `yaml:"gdock.5"`
	Gdock_6           map[string]any `yaml:"gdock.6"`
	Gdock_7           map[string]any `yaml:"gdock.7"`
	Gdock_8           map[string]any `yaml:"gdock.8"`
	Gdock_9           map[string]any `yaml:"gdock.9"`
	Gdock_10          map[string]any `yaml:"gdock.10"`
	Gdock_11          map[string]any `yaml:"gdock.11"`
	Gdock_12          map[string]any `yaml:"gdock.12"`
	Gdock_13          map[string]any `yaml:"gdock.13"`
	Gdock_14          map[string]any `yaml:"gdock.14"`
	Gdock_15          map[string]any `yaml:"gdock.15"`
	Gdock_16          map[string]any `yaml:"gdock.16"`
	Gdock_17          map[string]any `yaml:"gdock.17"`
	Gdock_18          map[string]any `yaml:"gdock.18"`
	Gdock_19          map[string]any `yaml:"gdock.19"`
	Gdock_20          map[string]any `yaml:"gdock.20"`
	Gdock_21          map[string]any `yaml:"gdock.21"`
	Gdock_22          map[string]any `yaml:"gdock.22"`
	Gdock_23          map[string]any `yaml:"gdock.23"`
	Gdock_24          map[string]any `yaml:"gdock.24"`
	Gdock_25          map[string]any `yaml:"gdock.25"`
	Lightdock         map[string]any `yaml:"lightdock"`
	Lightdock_0       map[string]any `yaml:"lightdock.0"`
	Lightdock_1       map[string]any `yaml:"lightdock.1"`
	Lightdock_2       map[string]any `yaml:"lightdock.2"`
	Lightdock_3       map[string]any `yaml:"lightdock.3"`
	Lightdock_4       map[string]any `yaml:"lightdock.4"`
	Lightdock_5       map[string]any `yaml:"lightdock.5"`
	Lightdock_6       map[string]any `yaml:"lightdock.6"`
	Lightdock_7       map[string]any `yaml:"lightdock.7"`
	Lightdock_8       map[string]any `yaml:"lightdock.8"`
	Lightdock_9       map[string]any `yaml:"lightdock.9"`
	Lightdock_10      map[string]any `yaml:"lightdock.10"`
	Lightdock_11      map[string]any `yaml:"lightdock.11"`
	Lightdock_12      map[string]any `yaml:"lightdock.12"`
	Lightdock_13      map[string]any `yaml:"lightdock.13"`
	Lightdock_14      map[string]any `yaml:"lightdock.14"`
	Lightdock_15      map[string]any `yaml:"lightdock.15"`
	Lightdock_16      map[string]any `yaml:"lightdock.16"`
	Lightdock_17      map[string]any `yaml:"lightdock.17"`
	Lightdock_18      map[string]any `yaml:"lightdock.18"`
	Lightdock_19      map[string]any `yaml:"lightdock.19"`
	Lightdock_20      map[string]any `yaml:"lightdock.20"`
	Lightdock_21      map[string]any `yaml:"lightdock.21"`
	Lightdock_22      map[string]any `yaml:"lightdock.22"`
	Lightdock_23      map[string]any `yaml:"lightdock.23"`
	Lightdock_24      map[string]any `yaml:"lightdock.24"`
	Lightdock_25      map[string]any `yaml:"lightdock.25"`
	Rigidbody         map[string]any `yaml:"rigidbody"`
	Rigidbody_0       map[string]any `yaml:"rigidbody.0"`
	Rigidbody_1       map[string]any `yaml:"rigidbody.1"`
	Rigidbody_2       map[string]any `yaml:"rigidbody.2"`
	Rigidbody_3       map[string]any `yaml:"rigidbody.3"`
	Rigidbody_4       map[string]any `yaml:"rigidbody.4"`
	Rigidbody_5       map[string]any `yaml:"rigidbody.5"`
	Rigidbody_6       map[string]any `yaml:"rigidbody.6"`
	Rigidbody_7       map[string]any `yaml:"rigidbody.7"`
	Rigidbody_8       map[string]any `yaml:"rigidbody.8"`
	Rigidbody_9       map[string]any `yaml:"rigidbody.9"`
	Rigidbody_10      map[string]any `yaml:"rigidbody.10"`
	Rigidbody_11      map[string]any `yaml:"rigidbody.11"`
	Rigidbody_12      map[string]any `yaml:"rigidbody.12"`
	Rigidbody_13      map[string]any `yaml:"rigidbody.13"`
	Rigidbody_14      map[string]any `yaml:"rigidbody.14"`
	Rigidbody_15      map[string]any `yaml:"rigidbody.15"`
	Rigidbody_16      map[string]any `yaml:"rigidbody.16"`
	Rigidbody_17      map[string]any `yaml:"rigidbody.17"`
	Rigidbody_18      map[string]any `yaml:"rigidbody.18"`
	Rigidbody_19      map[string]any `yaml:"rigidbody.19"`
	Rigidbody_20      map[string]any `yaml:"rigidbody.20"`
	Rigidbody_21      map[string]any `yaml:"rigidbody.21"`
	Rigidbody_22      map[string]any `yaml:"rigidbody.22"`
	Rigidbody_23      map[string]any `yaml:"rigidbody.23"`
	Rigidbody_24      map[string]any `yaml:"rigidbody.24"`
	Rigidbody_25      map[string]any `yaml:"rigidbody.25"`
	Emscoring         map[string]any `yaml:"emscoring"`
	Emscoring_0       map[string]any `yaml:"emscoring.0"`
	Emscoring_1       map[string]any `yaml:"emscoring.1"`
	Emscoring_2       map[string]any `yaml:"emscoring.2"`
	Emscoring_3       map[string]any `yaml:"emscoring.3"`
	Emscoring_4       map[string]any `yaml:"emscoring.4"`
	Emscoring_5       map[string]any `yaml:"emscoring.5"`
	Emscoring_6       map[string]any `yaml:"emscoring.6"`
	Emscoring_7       map[string]any `yaml:"emscoring.7"`
	Emscoring_8       map[string]any `yaml:"emscoring.8"`
	Emscoring_9       map[string]any `yaml:"emscoring.9"`
	Emscoring_10      map[string]any `yaml:"emscoring.10"`
	Emscoring_11      map[string]any `yaml:"emscoring.11"`
	Emscoring_12      map[string]any `yaml:"emscoring.12"`
	Emscoring_13      map[string]any `yaml:"emscoring.13"`
	Emscoring_14      map[string]any `yaml:"emscoring.14"`
	Emscoring_15      map[string]any `yaml:"emscoring.15"`
	Emscoring_16      map[string]any `yaml:"emscoring.16"`
	Emscoring_17      map[string]any `yaml:"emscoring.17"`
	Emscoring_18      map[string]any `yaml:"emscoring.18"`
	Emscoring_19      map[string]any `yaml:"emscoring.19"`
	Emscoring_20      map[string]any `yaml:"emscoring.20"`
	Emscoring_21      map[string]any `yaml:"emscoring.21"`
	Emscoring_22      map[string]any `yaml:"emscoring.22"`
	Emscoring_23      map[string]any `yaml:"emscoring.23"`
	Emscoring_24      map[string]any `yaml:"emscoring.24"`
	Emscoring_25      map[string]any `yaml:"emscoring.25"`
	Mdscoring         map[string]any `yaml:"mdscoring"`
	Mdscoring_0       map[string]any `yaml:"mdscoring.0"`
	Mdscoring_1       map[string]any `yaml:"mdscoring.1"`
	Mdscoring_2       map[string]any `yaml:"mdscoring.2"`
	Mdscoring_3       map[string]any `yaml:"mdscoring.3"`
	Mdscoring_4       map[string]any `yaml:"mdscoring.4"`
	Mdscoring_5       map[string]any `yaml:"mdscoring.5"`
	Mdscoring_6       map[string]any `yaml:"mdscoring.6"`
	Mdscoring_7       map[string]any `yaml:"mdscoring.7"`
	Mdscoring_8       map[string]any `yaml:"mdscoring.8"`
	Mdscoring_9       map[string]any `yaml:"mdscoring.9"`
	Mdscoring_10      map[string]any `yaml:"mdscoring.10"`
	Mdscoring_11      map[string]any `yaml:"mdscoring.11"`
	Mdscoring_12      map[string]any `yaml:"mdscoring.12"`
	Mdscoring_13      map[string]any `yaml:"mdscoring.13"`
	Mdscoring_14      map[string]any `yaml:"mdscoring.14"`
	Mdscoring_15      map[string]any `yaml:"mdscoring.15"`
	Mdscoring_16      map[string]any `yaml:"mdscoring.16"`
	Mdscoring_17      map[string]any `yaml:"mdscoring.17"`
	Mdscoring_18      map[string]any `yaml:"mdscoring.18"`
	Mdscoring_19      map[string]any `yaml:"mdscoring.19"`
	Mdscoring_20      map[string]any `yaml:"mdscoring.20"`
	Mdscoring_21      map[string]any `yaml:"mdscoring.21"`
	Mdscoring_22      map[string]any `yaml:"mdscoring.22"`
	Mdscoring_23      map[string]any `yaml:"mdscoring.23"`
	Mdscoring_24      map[string]any `yaml:"mdscoring.24"`
	Mdscoring_25      map[string]any `yaml:"mdscoring.25"`
	ProdigyLigand     map[string]any `yaml:"prodigyligand"`
	ProdigyLigand_0   map[string]any `yaml:"prodigyligand.0"`
	ProdigyLigand_1   map[string]any `yaml:"prodigyligand.1"`
	ProdigyLigand_2   map[string]any `yaml:"prodigyligand.2"`
	ProdigyLigand_3   map[string]any `yaml:"prodigyligand.3"`
	ProdigyLigand_4   map[string]any `yaml:"prodigyligand.4"`
	ProdigyLigand_5   map[string]any `yaml:"prodigyligand.5"`
	ProdigyLigand_6   map[string]any `yaml:"prodigyligand.6"`
	ProdigyLigand_7   map[string]any `yaml:"prodigyligand.7"`
	ProdigyLigand_8   map[string]any `yaml:"prodigyligand.8"`
	ProdigyLigand_9   map[string]any `yaml:"prodigyligand.9"`
	ProdigyLigand_10  map[string]any `yaml:"prodigyligand.10"`
	ProdigyLigand_11  map[string]any `yaml:"prodigyligand.11"`
	ProdigyLigand_12  map[string]any `yaml:"prodigyligand.12"`
	ProdigyLigand_13  map[string]any `yaml:"prodigyligand.13"`
	ProdigyLigand_14  map[string]any `yaml:"prodigyligand.14"`
	ProdigyLigand_15  map[string]any `yaml:"prodigyligand.15"`
	ProdigyLigand_16  map[string]any `yaml:"prodigyligand.16"`
	ProdigyLigand_17  map[string]any `yaml:"prodigyligand.17"`
	ProdigyLigand_18  map[string]any `yaml:"prodigyligand.18"`
	ProdigyLigand_19  map[string]any `yaml:"prodigyligand.19"`
	ProdigyLigand_20  map[string]any `yaml:"prodigyligand.20"`
	ProdigyLigand_21  map[string]any `yaml:"prodigyligand.21"`
	ProdigyLigand_22  map[string]any `yaml:"prodigyligand.22"`
	ProdigyLigand_23  map[string]any `yaml:"prodigyligand.23"`
	ProdigyLigand_24  map[string]any `yaml:"prodigyligand.24"`
	ProdigyLigand_25  map[string]any `yaml:"prodigyligand.25"`
	Prodigyprotein    map[string]any `yaml:"prodigyprotein"`
	Prodigyprotein_0  map[string]any `yaml:"prodigyprotein.0"`
	Prodigyprotein_1  map[string]any `yaml:"prodigyprotein.1"`
	Prodigyprotein_2  map[string]any `yaml:"prodigyprotein.2"`
	Prodigyprotein_3  map[string]any `yaml:"prodigyprotein.3"`
	Prodigyprotein_4  map[string]any `yaml:"prodigyprotein.4"`
	Prodigyprotein_5  map[string]any `yaml:"prodigyprotein.5"`
	Prodigyprotein_6  map[string]any `yaml:"prodigyprotein.6"`
	Prodigyprotein_7  map[string]any `yaml:"prodigyprotein.7"`
	Prodigyprotein_8  map[string]any `yaml:"prodigyprotein.8"`
	Prodigyprotein_9  map[string]any `yaml:"prodigyprotein.9"`
	Prodigyprotein_10 map[string]any `yaml:"prodigyprotein.10"`
	Prodigyprotein_11 map[string]any `yaml:"prodigyprotein.11"`
	Prodigyprotein_12 map[string]any `yaml:"prodigyprotein.12"`
	Prodigyprotein_13 map[string]any `yaml:"prodigyprotein.13"`
	Prodigyprotein_14 map[string]any `yaml:"prodigyprotein.14"`
	Prodigyprotein_15 map[string]any `yaml:"prodigyprotein.15"`
	Prodigyprotein_16 map[string]any `yaml:"prodigyprotein.16"`
	Prodigyprotein_17 map[string]any `yaml:"prodigyprotein.17"`
	Prodigyprotein_18 map[string]any `yaml:"prodigyprotein.18"`
	Prodigyprotein_19 map[string]any `yaml:"prodigyprotein.19"`
	Prodigyprotein_20 map[string]any `yaml:"prodigyprotein.20"`
	Prodigyprotein_21 map[string]any `yaml:"prodigyprotein.21"`
	Prodigyprotein_22 map[string]any `yaml:"prodigyprotein.22"`
	Prodigyprotein_23 map[string]any `yaml:"prodigyprotein.23"`
	Prodigyprotein_24 map[string]any `yaml:"prodigyprotein.24"`
	Prodigyprotein_25 map[string]any `yaml:"prodigyprotein.25"`
	Sasascore         map[string]any `yaml:"sasascore"`
	Sasascore_0       map[string]any `yaml:"sasascore.0"`
	Sasascore_1       map[string]any `yaml:"sasascore.1"`
	Sasascore_2       map[string]any `yaml:"sasascore.2"`
	Sasascore_3       map[string]any `yaml:"sasascore.3"`
	Sasascore_4       map[string]any `yaml:"sasascore.4"`
	Sasascore_5       map[string]any `yaml:"sasascore.5"`
	Sasascore_6       map[string]any `yaml:"sasascore.6"`
	Sasascore_7       map[string]any `yaml:"sasascore.7"`
	Sasascore_8       map[string]any `yaml:"sasascore.8"`
	Sasascore_9       map[string]any `yaml:"sasascore.9"`
	Sasascore_10      map[string]any `yaml:"sasascore.10"`
	Sasascore_11      map[string]any `yaml:"sasascore.11"`
	Sasascore_12      map[string]any `yaml:"sasascore.12"`
	Sasascore_13      map[string]any `yaml:"sasascore.13"`
	Sasascore_14      map[string]any `yaml:"sasascore.14"`
	Sasascore_15      map[string]any `yaml:"sasascore.15"`
	Sasascore_16      map[string]any `yaml:"sasascore.16"`
	Sasascore_17      map[string]any `yaml:"sasascore.17"`
	Sasascore_18      map[string]any `yaml:"sasascore.18"`
	Sasascore_19      map[string]any `yaml:"sasascore.19"`
	Sasascore_20      map[string]any `yaml:"sasascore.20"`
	Sasascore_21      map[string]any `yaml:"sasascore.21"`
	Sasascore_22      map[string]any `yaml:"sasascore.22"`
	Sasascore_23      map[string]any `yaml:"sasascore.23"`
	Sasascore_24      map[string]any `yaml:"sasascore.24"`
	Sasascore_25      map[string]any `yaml:"sasascore.25"`
	Topoaa            map[string]any `yaml:"topoaa"`
	Topoaa_0          map[string]any `yaml:"topoaa.0"`
	Topoaa_1          map[string]any `yaml:"topoaa.1"`
	Topoaa_2          map[string]any `yaml:"topoaa.2"`
	Topoaa_3          map[string]any `yaml:"topoaa.3"`
	Topoaa_4          map[string]any `yaml:"topoaa.4"`
	Topoaa_5          map[string]any `yaml:"topoaa.5"`
	Topoaa_6          map[string]any `yaml:"topoaa.6"`
	Topoaa_7          map[string]any `yaml:"topoaa.7"`
	Topoaa_8          map[string]any `yaml:"topoaa.8"`
	Topoaa_9          map[string]any `yaml:"topoaa.9"`
	Topoaa_10         map[string]any `yaml:"topoaa.10"`
	Topoaa_11         map[string]any `yaml:"topoaa.11"`
	Topoaa_12         map[string]any `yaml:"topoaa.12"`
	Topoaa_13         map[string]any `yaml:"topoaa.13"`
	Topoaa_14         map[string]any `yaml:"topoaa.14"`
	Topoaa_15         map[string]any `yaml:"topoaa.15"`
	Topoaa_16         map[string]any `yaml:"topoaa.16"`
	Topoaa_17         map[string]any `yaml:"topoaa.17"`
	Topoaa_18         map[string]any `yaml:"topoaa.18"`
	Topoaa_19         map[string]any `yaml:"topoaa.19"`
	Topoaa_20         map[string]any `yaml:"topoaa.20"`
	Topoaa_21         map[string]any `yaml:"topoaa.21"`
	Topoaa_22         map[string]any `yaml:"topoaa.22"`
	Topoaa_23         map[string]any `yaml:"topoaa.23"`
	Topoaa_24         map[string]any `yaml:"topoaa.24"`
	Topoaa_25         map[string]any `yaml:"topoaa.25"`
	Topocg            map[string]any `yaml:"topocg"`
	Topocg_0          map[string]any `yaml:"topocg.0"`
	Topocg_1          map[string]any `yaml:"topocg.1"`
	Topocg_2          map[string]any `yaml:"topocg.2"`
	Topocg_3          map[string]any `yaml:"topocg.3"`
	Topocg_4          map[string]any `yaml:"topocg.4"`
	Topocg_5          map[string]any `yaml:"topocg.5"`
	Topocg_6          map[string]any `yaml:"topocg.6"`
	Topocg_7          map[string]any `yaml:"topocg.7"`
	Topocg_8          map[string]any `yaml:"topocg.8"`
	Topocg_9          map[string]any `yaml:"topocg.9"`
	Topocg_10         map[string]any `yaml:"topocg.10"`
	Topocg_11         map[string]any `yaml:"topocg.11"`
	Topocg_12         map[string]any `yaml:"topocg.12"`
	Topocg_13         map[string]any `yaml:"topocg.13"`
	Topocg_14         map[string]any `yaml:"topocg.14"`
	Topocg_15         map[string]any `yaml:"topocg.15"`
	Topocg_16         map[string]any `yaml:"topocg.16"`
	Topocg_17         map[string]any `yaml:"topocg.17"`
	Topocg_18         map[string]any `yaml:"topocg.18"`
	Topocg_19         map[string]any `yaml:"topocg.19"`
	Topocg_20         map[string]any `yaml:"topocg.20"`
	Topocg_21         map[string]any `yaml:"topocg.21"`
	Topocg_22         map[string]any `yaml:"topocg.22"`
	Topocg_23         map[string]any `yaml:"topocg.23"`
	Topocg_24         map[string]any `yaml:"topocg.24"`
	Topocg_25         map[string]any `yaml:"topocg.25"`
	Exit              map[string]any `yaml:"exit"`
	Exit_0            map[string]any `yaml:"exit.0"`
	Exit_1            map[string]any `yaml:"exit.1"`
	Exit_2            map[string]any `yaml:"exit.2"`
	Exit_3            map[string]any `yaml:"exit.3"`
	Exit_4            map[string]any `yaml:"exit.4"`
	Exit_5            map[string]any `yaml:"exit.5"`
	Exit_6            map[string]any `yaml:"exit.6"`
	Exit_7            map[string]any `yaml:"exit.7"`
	Exit_8            map[string]any `yaml:"exit.8"`
	Exit_9            map[string]any `yaml:"exit.9"`
	Exit_10           map[string]any `yaml:"exit.10"`
	Exit_11           map[string]any `yaml:"exit.11"`
	Exit_12           map[string]any `yaml:"exit.12"`
	Exit_13           map[string]any `yaml:"exit.13"`
	Exit_14           map[string]any `yaml:"exit.14"`
	Exit_15           map[string]any `yaml:"exit.15"`
	Exit_16           map[string]any `yaml:"exit.16"`
	Exit_17           map[string]any `yaml:"exit.17"`
	Exit_18           map[string]any `yaml:"exit.18"`
	Exit_19           map[string]any `yaml:"exit.19"`
	Exit_20           map[string]any `yaml:"exit.20"`
	Exit_21           map[string]any `yaml:"exit.21"`
	Exit_22           map[string]any `yaml:"exit.22"`
	Exit_23           map[string]any `yaml:"exit.23"`
	Exit_24           map[string]any `yaml:"exit.24"`
	Exit_25           map[string]any `yaml:"exit.25"`
}
