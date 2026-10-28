import { StyleSheet } from "react-native";

export default StyleSheet.create({
  safeArea: {
    backgroundColor: "#0A5133",
    flex: 1,
    alignItems: "center",
  },

  // --- Cabeçalho ---
  containerCabecalho: {
    flexDirection: "row",
    alignItems: "center",
    marginTop: 20,
    paddingHorizontal: 20,
    width: "100%",
  },
  tituloCabecalho: {
    flex: 1,
    textAlign: "center",
    fontSize: 22,
    fontWeight: "bold",
    color: "#FFFFFF",
    marginRight: 28,
  },

  // --- Cartão Principal ---
  suporteCartao: {
    width: "90%",
    height: "75%",
    alignSelf: "center",
    marginTop: 20,
  },
  containerPrincipal: {
    width: "100%",
    height: "100%",
    backgroundColor: "#D9D9D9",
    borderRadius: 16,
    borderWidth: 1.5,
    borderColor: "#000000",
    overflow: "hidden",
    shadowColor: '#000', 
    shadowOffset: { width: 0, height: 2 }, 
    shadowOpacity: 0.1, 
    shadowRadius: 4, 
    elevation: 3 
  },

  // --- Bloco Superior Interno ---
  blocoSuperior: {
    flexDirection: "row",
    padding: 15,
  },
  imagemAnimal: {
    width: 100,
    height: 100,
    borderRadius: 8,
    borderWidth: 1,
    borderColor: "#A0A0A0",
  },
  divisorVertical: {
    width: 2,
    backgroundColor: "#000000",
    marginHorizontal: 12,
  },
  containerInformacoes: {
    flex: 1,
  },
  nomeAnimal: {
    fontWeight: "bold",
    textTransform: "uppercase",
    fontSize: 14,
    color: "#000000",
  },
  divisorHorizontal: {
    height: 2,
    backgroundColor: "#000000",
    marginVertical: 4,
  },
  nomeCientifico: {
    fontStyle: "italic",
    fontSize: 12,
    marginBottom: 6,
    color: "#222222",
  },
  textoInformacao: {
    fontSize: 12,
    fontWeight: "600",
    color: "#000000",
    marginTop: 2,
  },
  blocoDescricao: {
    padding: 15,
    borderTopWidth: 1,
    borderTopColor: "#A0A0A0",
  },
  tituloDescricao: {
    fontSize: 15,
    fontWeight: "bold",
    color: "#000000",
    marginBottom: 8,
  },
  tituloDescricao2: {
    fontSize: 15,
    fontWeight: "bold",
    color: "#333333",
    marginBottom: 8,
  },
  textoDescricao: {
    fontSize: 14,
    color: "#333333",
    lineHeight: 20,
  },
  linhaInformacao: {
  flexDirection: "row",
  justifyContent: "space-between",
  alignItems: "center",
  marginTop: 4,
  },
  valorInformacao: {
    fontSize: 12,
    fontWeight: "bold",
    color: "#333333",
  },
  linhaInformacao2: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
    marginTop: 4,
  },
  linhaInformacao3: {
    flexDirection: "row",
    alignItems: "center",
    marginTop: 4,
  },
  textoInformacao2: {
    fontSize: 14,
    fontWeight: "600",
    color: "#000000",
    marginTop: 2,
  },
  valorInformacao2: {
    fontSize: 14,
    color: "#333333",
    },
    containerAviso: {
    backgroundColor: 'rgba(255, 140, 0, 0.15)', // Laranja com 15% de opacidade
    borderRadius: 8,
    padding: 12,
    marginTop: 16,
    borderWidth: 1,
    borderColor: 'rgba(255, 140, 0, 0.4)', // Borda sutil para definir melhor o quadrado
  },
  cabecalhoAviso: {
    flexDirection: 'row',
    alignItems: 'center',
    marginBottom: 6,
  },
  iconeAviso: {
    marginRight: 6,
  },
  tituloAviso: {
    fontWeight: 'bold',
    fontSize: 14,
    color: '#BF360C',
  },
  textoAviso: {
    fontSize: 13,
    color: '#BF360C',
    marginTop: 2,
  },
  modalContainer: {
    flex: 1,
    backgroundColor: "rgba(0, 0, 0, 0.9)",
    justifyContent: "center",
    alignItems: "center",
  },
  botaoFecharModal: {
    position: "absolute",
    top: 40,
    right: 20,
    zIndex: 10,
    padding: 10,
  },
  imagemModal: {
    width: "90%",
    height: "80%",
  },

});